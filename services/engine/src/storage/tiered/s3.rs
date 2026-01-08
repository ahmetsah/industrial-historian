#![allow(dead_code)]
use anyhow::{Context, Result};
use s3::bucket::Bucket;
use s3::creds::Credentials;
use s3::region::Region;

#[derive(Clone)]
pub struct S3Client {
    bucket: Box<Bucket>,
}

impl S3Client {
    pub fn new(
        endpoint: &str,
        bucket_name: &str,
        access_key: &str,
        secret_key: &str,
    ) -> Result<Self> {
        let region = Region::Custom {
            region: "us-east-1".to_owned(),
            endpoint: endpoint.to_owned(),
        };
        let credentials = Credentials::new(Some(access_key), Some(secret_key), None, None, None)?;
        let mut bucket = Bucket::new(bucket_name, region, credentials)?;
        bucket.set_path_style();
        Ok(Self {
            bucket: Box::new(bucket),
        })
    }

    pub async fn put_object(&self, key: &str, data: &[u8]) -> Result<()> {
        let response = self
            .bucket
            .put_object(key, data)
            .await
            .context("Failed to put object")?;

        // Basic integrity check: Verify status code is 200 OK
        if response.status_code() != 200 {
            anyhow::bail!("S3 upload failed with status: {}", response.status_code());
        }
        Ok(())
    }

    pub async fn get_object(&self, key: &str) -> Result<Vec<u8>> {
        let response = self
            .bucket
            .get_object(key)
            .await
            .context("Failed to get object")?;
        Ok(response.to_vec())
    }

    pub async fn delete_object(&self, key: &str) -> Result<()> {
        self.bucket
            .delete_object(key)
            .await
            .context("Failed to delete object")?;
        Ok(())
    }
}
