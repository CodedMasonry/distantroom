pub mod profile;
pub mod logging;

#[derive(Debug)]
pub struct State {
    pub profile: profile::Profile,
}
