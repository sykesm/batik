// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

mod messages;

use messages::{resolved, transaction, validation_api};
use p256::ecdsa::{Signature, VerifyingKey};
use protobuf::{parse_from_bytes, Message};
use signature::Verifier;
use simple_asn1::{oid, ASN1Block, ASN1DecodeErr, BigUint, OID};

#[link(wasm_import_module = "batik")]
extern "C" {
    #[link_name = "log"]
    fn __batik_log(msg: *const u8, len: usize);

    #[link_name = "read"]
    fn __batik_read(stream_id: isize, buf: *mut u8, count: usize) -> isize;

    #[link_name = "write"]
    fn __batik_write(stream_id: isize, buf: *const u8, count: usize) -> isize;
}

fn batik_log(msg: &str) {
    unsafe { __batik_log(msg.as_ptr(), msg.len()) }
}

fn batik_read(id: i32, buf: &mut Vec<u8>) -> isize {
    let len = unsafe { __batik_read(id as isize, buf.as_mut_ptr() as *mut u8, buf.capacity()) };
    if len >= 0 {
        unsafe { buf.set_len(len as usize) };
    }
    len
}

fn batik_write(id: i32, buf: &Vec<u8>) -> isize {
    unsafe { __batik_write(id as isize, buf.as_ptr() as *const u8, buf.len()) }
}

fn extract_sec1_key(pkix_key: &[u8]) -> Result<Vec<u8>, ASN1DecodeErr> {
    let der = simple_asn1::from_der(pkix_key)?;
    let der = der.first().unwrap();
    if let ASN1Block::Sequence(_, seq) = der {
        if seq.len() != 2 {
            return Err(ASN1DecodeErr::Incomplete);
        }
        if let ASN1Block::Sequence(_, algid) = seq.get(0).unwrap() {
            if let ASN1Block::ObjectIdentifier(_, alg) = algid.get(0).unwrap() {
                if alg != oid!(1, 2, 840, 10045, 2, 1) {
                    // ecPublicKey OID from RFC 3279 / RFC 5753
                    return Err(ASN1DecodeErr::Incomplete);
                }
            }
            if let ASN1Block::ObjectIdentifier(_, params) = algid.get(1).unwrap() {
                if params != oid!(1, 2, 840, 10045, 3, 1, 7) {
                    // P256 OID from RFC 5759
                    return Err(ASN1DecodeErr::Incomplete);
                }
            }
        }
        if let ASN1Block::BitString(_, _, pk) = seq.get(1).unwrap() {
            return Ok(pk.to_vec());
        }
    }
    Ok(pkix_key.to_vec())
}

#[no_mangle]
pub extern "C" fn validate(stream: i32, input_len: i32) -> i32 {
    let req_bytes: &mut Vec<u8> = &mut Vec::with_capacity(input_len as usize);
    let read_len = batik_read(stream, req_bytes);
    if read_len != input_len as isize {
        return -1;
    }

    // Decode the bytes into the ResolvedTransaction protobuf
    let request = parse_from_bytes::<validation_api::ValidateRequest>(&req_bytes).unwrap();
    let resolved_tx = request.get_resolved_transaction();
    let txid = resolved_tx.get_txid();
    batik_log(format!("txid {}", hex::encode(txid)).as_str());

    for signer in required_signers(resolved_tx) {
        let signatures = resolved_tx.get_signatures();
        let sig = signature(signer.get_public_key(), signatures.to_vec());
        batik_log(format!("sig {}", hex::encode(sig.get_signature())).as_str());

        let pkix = sig.get_public_key();
        batik_log(format!("pkix {}", hex::encode(pkix)).as_str());
        let sec1 = extract_sec1_key(pkix).unwrap();
        batik_log(format!("pk: {:?}", hex::encode(&sec1)).as_str());
        let verifying_key = VerifyingKey::from_sec1_bytes(&sec1).unwrap();
        batik_log(format!("verifying_key {:?}", verifying_key).as_str());
        let signature = Signature::from_asn1(sig.get_signature()).unwrap();

        if !verifying_key.verify(txid, &signature).is_ok() {
            batik_log(format!("signature is not valid!").as_str());
            return -1;
        }
    }

    let result = validation_api::ValidateResponse::new();
    let encoded_result = result.write_to_bytes().unwrap();
    if batik_write(stream, &encoded_result) != encoded_result.len() as isize {
        return -1;
    }
    0
}

fn required_signers(resolved: &resolved::ResolvedTransaction) -> Vec<transaction::Party> {
    let mut required = Vec::new();
    for input in resolved.get_inputs() {
        required.append(&mut input.get_state().get_info().get_owners().to_vec());
    }
    required.append(&mut resolved.get_required_signers().to_vec());
    required
}

fn signature(public_key: &[u8], signatures: Vec<transaction::Signature>) -> transaction::Signature {
    for sig in signatures {
        if sig.public_key == public_key {
            return sig;
        }
    }
    return transaction::Signature::new();
}
