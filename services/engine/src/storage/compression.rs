use anyhow::Result;
use tsz::stream::{BufferedReader, BufferedWriter};
use tsz::{DataPoint, Decode, Encode, StdDecoder, StdEncoder};

pub fn compress_points(points: &[(i64, f64)]) -> Result<Vec<u8>> {
    if points.is_empty() {
        return Ok(Vec::new());
    }
    let start = points[0].0 as u64;
    let writer = BufferedWriter::new();
    let mut encoder = StdEncoder::new(start, writer);
    for (t, v) in points {
        encoder.encode(DataPoint::new(*t as u64, *v));
    }
    Ok(encoder.close().into_vec())
}

pub fn decompress_points(data: &[u8]) -> Result<Vec<(i64, f64)>> {
    let reader = BufferedReader::new(data.to_vec().into_boxed_slice());
    let mut decoder = StdDecoder::new(reader);
    let mut points = Vec::new();

    while let Ok(dp) = decoder.next() {
        points.push((dp.get_time() as i64, dp.get_value()));
    }

    Ok(points)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_compression_decompression() {
        let points = vec![
            (1000, 1.0),
            (2000, 1.1),
            (3000, 1.2),
            (4000, 1.2), // Delta 0
            (5000, 1.2),
        ];

        let compressed = compress_points(&points).unwrap();
        assert!(!compressed.is_empty());
        // Compressed size should be smaller than raw (5 * 16 = 80 bytes)
        // With Gorilla, it should be significantly smaller for this regular data
        println!("Compressed size: {}", compressed.len());

        let decompressed = decompress_points(&compressed).unwrap();
        assert_eq!(points, decompressed);
    }
}
