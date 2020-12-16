// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

mod batik;
mod messages;

use messages::resolved::ResolvedTransaction;
use messages::transaction::{Party, Signature};
use messages::validation_api::{ValidateRequest, ValidateResponse};
use p256::ecdsa;
use protobuf::{parse_from_bytes, Message};
use signature::Verifier;
use simple_asn1::{oid, BigUint, OID};

#[derive(Debug)]
enum Error {
    InvalidAlgorithmEncoding,
    InvalidDer(simple_asn1::ASN1DecodeErr),
    InvalidKeyEncoding,
    InvalidPKIXEncoding,
    MissingSignature(Party),
    ProtobufError(protobuf::ProtobufError),
    UnknownAlgorithm,
    UnknownCurve,
    SignatureError(signature::Error),
}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::result::Result<(), std::fmt::Error> {
        match self {
            Error::InvalidAlgorithmEncoding => {
                f.write_str("invalid ASN.1 encoding for public key algorithm")?
            }
            Error::InvalidDer(e) => e.fmt(f)?,
            Error::InvalidKeyEncoding => f.write_str("invalid ASN.1 encoding for public key")?,
            Error::InvalidPKIXEncoding => {
                f.write_str("invalid ASN.1 encoding for subject public key")?
            }
            Error::MissingSignature(party) => {
                let pk = party.get_public_key();
                f.write_fmt(format_args!("missing signature for {}", hex::encode(pk)))
            }?,
            Error::ProtobufError(e) => e.fmt(f)?,
            Error::UnknownAlgorithm => f.write_str("unknown algorithm")?,
            Error::UnknownCurve => f.write_str("unknown curve")?,
            Error::SignatureError(e) => e.fmt(f)?,
        }
        Ok(())
    }
}

impl From<simple_asn1::ASN1DecodeErr> for Error {
    fn from(err: simple_asn1::ASN1DecodeErr) -> Error {
        Error::InvalidDer(err)
    }
}

impl From<protobuf::ProtobufError> for Error {
    fn from(err: protobuf::ProtobufError) -> Error {
        Error::ProtobufError(err)
    }
}

impl From<signature::Error> for Error {
    fn from(err: signature::Error) -> Error {
        Error::SignatureError(err)
    }
}

// Type alias that makes use of the local Error type.
type Result<T> = std::result::Result<T, Error>;

#[no_mangle]
pub extern "C" fn validate(stream: i32, input_len: i32) -> i32 {
    let req_bytes: &mut Vec<u8> = &mut Vec::with_capacity(input_len as usize);
    if batik::read(stream, req_bytes) != input_len as isize {
        return -1;
    }

    match validate_tx(req_bytes) {
        Ok(res) if batik::write(stream, &res) == res.len() as isize => 0,
        _ => -1,
    }
}

fn validate_tx(req_bytes: &mut Vec<u8>) -> Result<Vec<u8>> {
    let request: ValidateRequest = parse_from_bytes(req_bytes)?;
    let tx = request.get_resolved_transaction();

    let mut resp = ValidateResponse::new();
    match verify_sigs(&tx) {
        Ok(_) => resp.set_valid(true),
        Err(e) => resp.set_error_message(format!("{}", e)),
    }

    resp.write_to_bytes().map_err(|e| Error::ProtobufError(e))
}

fn verify_sigs(tx: &ResolvedTransaction) -> Result<()> {
    let txid = tx.get_txid();
    let signatures = tx.get_signatures();
    for signer in required_signers(tx) {
        let pkix_key = signer.get_public_key();

        let sec1 = extract_sec1_key(pkix_key)?;
        let vk = ecdsa::VerifyingKey::from_sec1_bytes(&sec1)?;

        let sig = signature(&signatures, pkix_key).ok_or(Error::MissingSignature(signer))?;
        let signature = ecdsa::Signature::from_asn1(sig.get_signature())?;

        vk.verify(txid, &signature)?;
    }
    Ok(())
}

fn required_signers(tx: &ResolvedTransaction) -> Vec<Party> {
    let mut required = tx.get_required_signers().to_vec();
    for input in tx.get_inputs() {
        required.append(&mut input.get_state().get_info().get_owners().to_vec());
    }
    required
}

fn signature<'a>(signatures: &'a [Signature], public_key: &[u8]) -> Option<&'a Signature> {
    signatures.iter().find(|sig| sig.public_key == public_key)
}

fn ec_public_key_oid() -> simple_asn1::OID {
    oid!(1, 2, 840, 10045, 2, 1)
}

fn ec_p256v1_oid() -> simple_asn1::OID {
    oid!(1, 2, 840, 10045, 3, 1, 7)
}

fn extract_sec1_key(pkix_subject_key: &[u8]) -> Result<Vec<u8>> {
    let der = simple_asn1::from_der(pkix_subject_key)?;
    let block = der.first().ok_or(Error::InvalidPKIXEncoding)?;
    let seq = match &block {
        simple_asn1::ASN1Block::Sequence(_, seq) if seq.len() == 2 => seq,
        _ => return Err(Error::InvalidPKIXEncoding),
    };
    let alg_id = match &seq[0] {
        simple_asn1::ASN1Block::Sequence(_, alg_id) if alg_id.len() == 2 => alg_id,
        _ => return Err(Error::InvalidAlgorithmEncoding),
    };
    let alg = match &alg_id[0] {
        simple_asn1::ASN1Block::ObjectIdentifier(_, alg) => alg,
        _ => return Err(Error::InvalidAlgorithmEncoding),
    };
    let curve = match &alg_id[1] {
        simple_asn1::ASN1Block::ObjectIdentifier(_, curve) => curve,
        _ => return Err(Error::InvalidAlgorithmEncoding),
    };
    if alg != ec_public_key_oid() {
        return Err(Error::UnknownAlgorithm);
    }
    if curve != ec_p256v1_oid() {
        return Err(Error::UnknownCurve);
    }
    let pk = match &seq[1] {
        simple_asn1::ASN1Block::BitString(_, _, pk) => pk,
        _ => return Err(Error::InvalidKeyEncoding),
    };

    Ok(pk.to_vec())
}

#[cfg(test)]
mod tests {
    use super::*;
    use simple_asn1::ASN1Block;

    macro_rules! assert_error_match {
        ($expression:expr, $error:pat) => {
            match $expression {
                Err($error) => (),
                e => assert!(false, "expected {:?}, got {:?}", stringify!($error), e),
            }
        };
    }

    #[test]
    fn sec1_key_extraction_empty_block() {
        let empty: Vec<u8> = Vec::new();
        assert_error_match!(extract_sec1_key(&empty), Error::InvalidDer(_));
    }

    #[test]
    fn sec1_key_extraction_not_sequence() {
        let block = ASN1Block::Boolean(0, true);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::InvalidPKIXEncoding);
    }

    #[test]
    fn sec1_key_extraction_not_sequence_len2() {
        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Boolean(0, true));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::InvalidPKIXEncoding);
    }

    #[test]
    fn sec1_key_extraction_algid_not_sequence() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::Boolean(0, true));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::InvalidAlgorithmEncoding);
    }

    #[test]
    fn sec1_key_extraction_bad_algid_element0() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::Boolean(0, true));
        algid.push(ASN1Block::Boolean(0, true));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::InvalidAlgorithmEncoding);
    }

    #[test]
    fn sec1_key_extraction_bad_algid_element1() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));
        algid.push(ASN1Block::Boolean(0, true));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::InvalidAlgorithmEncoding);
    }

    #[test]
    fn sec1_key_extraction_bad_algid_unknown_algorithm() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));
        algid.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::UnknownAlgorithm);
    }

    #[test]
    fn sec1_key_extraction_bad_algid_unknown_curve() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::ObjectIdentifier(0, ec_public_key_oid()));
        algid.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::UnknownCurve);
    }

    #[test]
    fn sec1_key_extraction_invalid_key_encoding() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::ObjectIdentifier(0, ec_public_key_oid()));
        algid.push(ASN1Block::ObjectIdentifier(0, ec_p256v1_oid()));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 3)));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_error_match!(extract_sec1_key(&pkix_key), Error::InvalidKeyEncoding);
    }

    #[test]
    fn sec1_key_extraction_happy() {
        let mut algid: Vec<ASN1Block> = Vec::new();
        algid.push(ASN1Block::ObjectIdentifier(0, oid!(1, 2, 840, 10045, 2, 1)));
        algid.push(ASN1Block::ObjectIdentifier(0, ec_p256v1_oid()));

        let mut seq: Vec<ASN1Block> = Vec::new();
        seq.push(ASN1Block::Sequence(0, algid));
        seq.push(ASN1Block::BitString(0, 1, vec![1u8, 2, 3]));

        let block = ASN1Block::Sequence(0, seq);
        let pkix_key = simple_asn1::to_der(&block).unwrap();
        assert_eq!(extract_sec1_key(&pkix_key).unwrap(), vec![1u8, 2, 3]);
    }
}
