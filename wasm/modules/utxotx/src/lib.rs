// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

mod batik;
mod messages;

use messages::{resolved, transaction, validation_api};
use p256::ecdsa::{Signature, VerifyingKey};
use protobuf::{parse_from_bytes, Message};
use signature::Verifier;
use simple_asn1::{oid, ASN1Block, ASN1DecodeErr, BigUint, OID};

#[no_mangle]
pub extern "C" fn validate(stream: i32, input_len: i32) -> i32 {
    let req_bytes: &mut Vec<u8> = &mut Vec::with_capacity(input_len as usize);
    if batik::read(stream, req_bytes) != input_len as isize {
        return -1;
    }

    match validate_tx(req_bytes) {
        Ok(res) if batik::write(stream, &res) == res.len() as isize => 0,
        Err(_) => -1,
        _ => -1,
    }
}

fn validate_tx(req_bytes: &mut Vec<u8>) -> Result<Vec<u8>, ASN1DecodeErr> {
    // Decode the bytes into the ResolvedTransaction protobuf
    let request = parse_from_bytes::<validation_api::ValidateRequest>(&req_bytes).unwrap();
    let resolved_tx = request.get_resolved_transaction();
    let txid = resolved_tx.get_txid();
    batik::log(format!("txid {}", hex::encode(txid)).as_str());

    let signatures = resolved_tx.get_signatures().to_vec();
    for signer in required_signers(resolved_tx) {
        let sig = signature(&signatures, signer.get_public_key()).unwrap();
        batik::log(format!("sig {}", hex::encode(sig.get_signature())).as_str());

        let pkix = sig.get_public_key();
        batik::log(format!("pkix {}", hex::encode(pkix)).as_str());
        let sec1 = extract_sec1_key(pkix).unwrap();
        batik::log(format!("pk {}", hex::encode(&sec1)).as_str());
        let vk = VerifyingKey::from_sec1_bytes(&sec1).unwrap();
        batik::log(format!("vk {}", hex::encode(vk.to_encoded_point(false).to_bytes()),).as_str());
        let signature = Signature::from_asn1(sig.get_signature()).unwrap();

        if !vk.verify(txid, &signature).is_ok() {
            batik::log(format!("signature is not valid!").as_str());
            return Err(ASN1DecodeErr::Incomplete);
        }
    }

    validation_api::ValidateResponse::new()
        .write_to_bytes()
        .map_err(|_| ASN1DecodeErr::Incomplete)
}

fn required_signers(resolved: &resolved::ResolvedTransaction) -> Vec<transaction::Party> {
    let mut required = resolved.get_required_signers().to_vec();
    for input in resolved.get_inputs() {
        required.append(&mut input.get_state().get_info().get_owners().to_vec());
    }
    required
}

fn signature<'a>(
    signatures: &'a Vec<transaction::Signature>,
    public_key: &[u8],
) -> Option<&'a transaction::Signature> {
    signatures.iter().find(|sig| sig.public_key == public_key)
}

fn extract_sec1_key(pkix_subject_key: &[u8]) -> Result<Vec<u8>, ASN1DecodeErr> {
    let der = simple_asn1::from_der(pkix_subject_key)?; // map_err ?
    let block = der.first().ok_or_else(|| ASN1DecodeErr::Incomplete)?;
    let seq = match &block {
        ASN1Block::Sequence(_, seq) if seq.len() == 2 => seq,
        _ => return Err(ASN1DecodeErr::Incomplete),
    };
    let alg_id = match &seq[0] {
        ASN1Block::Sequence(_, alg_id) if alg_id.len() == 2 => alg_id,
        _ => return Err(ASN1DecodeErr::Incomplete),
    };
    let alg = match &alg_id[0] {
        ASN1Block::ObjectIdentifier(_, alg) => alg,
        _ => return Err(ASN1DecodeErr::Incomplete),
    };
    let curve = match &alg_id[1] {
        ASN1Block::ObjectIdentifier(_, curve) => curve,
        _ => return Err(ASN1DecodeErr::Incomplete),
    };
    if alg != oid!(1, 2, 840, 10045, 2, 1) {
        return Err(ASN1DecodeErr::Incomplete);
    }
    if curve != oid!(1, 2, 840, 10045, 3, 1, 7) {
        return Err(ASN1DecodeErr::Incomplete);
    }
    let pk = match &seq[1] {
        ASN1Block::BitString(_, _, pk) => pk,
        _ => return Err(ASN1DecodeErr::Incomplete),
    };

    Ok(pk.to_vec())
}
