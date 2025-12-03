//use historian_core::hello;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    tracing::info!("Starting Ingestor Service");
    //hello();
}
