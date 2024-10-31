use std::{thread, time::Duration};

use futures_util::{SinkExt, StreamExt, TryStreamExt};
use reqwest::{Certificate, Identity};
use reqwest_websocket::{Message, RequestBuilderExt};

use crate::Profile;

#[tokio::main(flavor = "current_thread")]
pub async fn connect(profile: &Profile) -> Result<(), anyhow::Error> {
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

    let response: reqwest_websocket::WebSocket = client
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

    let (mut tx, mut rx) = response.split();
    tokio::task::spawn_local(async move {
        for i in 1..11 {
            tx.send(Message::Text(format!("Hello, World! #{i}")))
                .await
                .unwrap();
        }
    });

    while let Some(message) = rx.try_next().await? {
        if let Message::Text(text) = message {
            println!("received: {text}");
        }
    }

    thread::sleep(Duration::from_secs(5));
    Ok(())
}
