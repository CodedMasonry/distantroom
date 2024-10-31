use log::debug;
use reqwest::{Certificate, Identity};

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
        .use_rustls_tls()
        .build()?;

    debug!("Connecting to Server...");
    let response = client
        .get(format!(
            "https://{}:{}/status",
            profile.inner.host.clone(),
            profile.inner.port
        ))
        .send()
        .await?;

    debug!("Response Received!");
    debug!("{:#?}", response);

    Ok(())
}
