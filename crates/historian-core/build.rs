fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(true)
        .build_client(true)
        .type_attribute(".", "#[derive(serde::Serialize, serde::Deserialize)]") // Add Serde support
        .compile_protos(
            &[
                "src/proto/common.proto",
                "src/proto/query.proto",
                "src/proto/analytics.proto",
            ],
            &["src/proto"],
        )?;
    Ok(())
}
