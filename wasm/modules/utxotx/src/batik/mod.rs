// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

#[link(wasm_import_module = "batik")]
extern "C" {
    #[link_name = "log"]
    fn __batik_log(msg: *const u8, len: usize);

    #[link_name = "read"]
    fn __batik_read(stream_id: isize, buf: *mut u8, count: usize) -> isize;

    #[link_name = "write"]
    fn __batik_write(stream_id: isize, buf: *const u8, count: usize) -> isize;
}

#[allow(dead_code)]
pub fn log(msg: &str) {
    unsafe { __batik_log(msg.as_ptr(), msg.len()) }
}

pub fn read(id: i32, buf: &mut Vec<u8>) -> isize {
    let len = unsafe { __batik_read(id as isize, buf.as_mut_ptr() as *mut u8, buf.capacity()) };
    if len >= 0 {
        unsafe { buf.set_len(len as usize) };
    }
    len
}

pub fn write(id: i32, buf: &Vec<u8>) -> isize {
    unsafe { __batik_write(id as isize, buf.as_ptr() as *const u8, buf.len()) }
}
