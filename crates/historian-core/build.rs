fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(false)
        .build_client(false) // We only need types for now, client/server code might be needed later but story says "shared library with Protobuf definitions"
        .type_attribute(".", "#[derive(serde::Serialize, serde::Deserialize)]") // Add Serde support
        .compile_protos(&["src/proto/common.proto"], &["src/proto"])?;
    Ok(())
}
