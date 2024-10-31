use futures_util::{
    stream::{SplitSink, SplitStream},
    StreamExt,
};
use reqwest::{Certificate, Identity};
use reqwest_websocket::{Message, RequestBuilderExt, WebSocket};
use tokio::runtime::Runtime;

use crate::Profile;

pub struct Server {
    tx: SplitSink<WebSocket, Message>,
    rx: SplitStream<WebSocket>,
    /// A `current_thread` runtime for executing operations on the
    /// asynchronous client in a blocking manner.
    rt: Runtime,
}
impl Server {
    pub fn connect(profile: &Profile) -> Result<Server, anyhow::Error> {
        let rt: Runtime = tokio::runtime::Builder::new_current_thread()
            .enable_all()
            .build()?;

        let response = rt.block_on(establish_conn(profile))?;

        let (tx, rx) = response.split();
        let server = Server { tx, rx, rt };

        Ok(server)
    }
}

async fn establish_conn(profile: &Profile) -> Result<WebSocket, anyhow::Error> {
    let format_pem = format!(
        "{}\n\n{}",
        profile.get_client_key(),
        profile.get_client_cert()
    );
    let root_cert = Certificate::from_pem(profile.get_root_cert().as_bytes())?;
    let identity = Identity::from_pem(format_pem.as_bytes())?;

    let client = reqwest::ClientBuilder::new()
        .add_root_certificate(root_cert)
        .identity(identity)
        .build()?;

    let socket = client
        .get(format!(
            "wss://{}:{}/ws/10",
            profile.inner.host.clone(),
            profile.inner.port
        ))
        .upgrade()
        .send()
        .await?
        .into_websocket()
        .await?;

    Ok(socket)
}
