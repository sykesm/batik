// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

use protobuf_codegen_pure;
use std::path::Path;

fn main() {
    let protos = &[
        "../../../protos/tx/v1/transaction.proto",
        "../../../protos/validation/v1/validation_api.proto",
        "../../../protos/validation/v1/resolved.proto",
    ];
    for proto in protos {
        println!("cargo:rerun-if-changed={}", proto);
    }
    generate_pb_rs(protos);
}

fn generate_pb_rs(protos: impl IntoIterator<Item = impl AsRef<Path>>) {
    protobuf_codegen_pure::Codegen::new()
        .out_dir("src/messages")
        .inputs(protos)
        .includes(&["../../../protos"])
        .run()
        .expect("Running protoc failed.");
}
