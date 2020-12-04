// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

extern crate protobuf_codegen_pure;

fn main() {
    println!("cargo:rerun-if-changed=../../../protos");
    generate_pb_rs();
}

fn generate_pb_rs() {
    protobuf_codegen_pure::Codegen::new()
        .out_dir("src/messages")
        .inputs(&[
            "../../../protos/tx/v1/transaction.proto",
            "../../../protos/validation/v1/validation_api.proto",
            "../../../protos/validation/v1/resolved.proto",
        ])
        .includes(&["../../../protos"])
        .run()
        .expect("Running protoc failed.");
}
