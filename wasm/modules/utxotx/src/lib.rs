// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

mod messages;

use messages::validation_api;
use protobuf::{parse_from_bytes, Message};
use std::os::raw::c_void;

#[link(wasm_import_module = "batik")]
extern "C" {
    #[link_name = "log"]
    fn __batik_log(msg: *const u8, len: usize);

    #[link_name = "read"]
    fn __batik_read(stream_id: isize, buf: *mut c_void, count: usize) -> isize;

    #[link_name = "write"]
    fn __batik_write(stream_id: isize, buf: *const c_void, count: usize) -> isize;
}

fn batik_log(msg: &str) {
    unsafe { __batik_log(msg.as_ptr(), msg.len()) }
}

fn batik_read(id: i32, buf: &mut Vec<u8>) -> isize {
    let len = unsafe { __batik_read(id as isize, buf.as_mut_ptr() as *mut c_void, buf.capacity()) };
    if len >= 0 {
        unsafe { buf.set_len(len as usize) };
    }
    len
}

fn batik_write(id: i32, buf: &Vec<u8>) -> isize {
    unsafe { __batik_write(id as isize, buf.as_ptr() as *const c_void, buf.len()) }
}

fn u32_from_slice(buf: &[u8]) -> u32 {
    let mut array = [0; 4];
    let b = &buf[..array.len()];
    array.copy_from_slice(b);
    u32::from_be_bytes(array)
}

#[no_mangle]
pub extern "C" fn validate(stream: i32) -> i32 {
    batik_log(format!("stream {}", stream).as_str());

    let len_buf: &mut Vec<u8> = &mut Vec::with_capacity(4);
    let read_len_rc = batik_read(stream, len_buf);
    batik_log(format!("stream {} read {}", stream, read_len_rc).as_str());
    if read_len_rc != 4 {
        return -1;
    }

    let input_len = u32_from_slice(&len_buf[0..4]);
    let req_bytes: &mut Vec<u8> = &mut Vec::with_capacity(input_len as usize);
    let read_len = batik_read(stream, req_bytes);
    if read_len != input_len as isize {
        return -1;
    }

    // Decode the bytes into the ResolvedTransaction protobuf
    let request = parse_from_bytes::<validation_api::ValidateRequest>(&req_bytes).unwrap();
    batik_log(format!("txid {:?}", request.get_resolved_transaction().get_txid()).as_str());

    // Create a ValidationResponse
    // let mut result = validation_api::ValidateResponse::new();
    let result = validation_api::ValidateResponse::new();

    let encoded_result = result.write_to_bytes().unwrap();
    let resp_len = u32::to_be_bytes(encoded_result.len() as u32).to_vec();
    if batik_write(stream, &resp_len) != 4 {
        return -1;
    }
    if batik_write(stream, &encoded_result) != encoded_result.len() as isize {
        return -1;
    }
    0
}
