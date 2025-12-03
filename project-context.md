# PROJE BAĞLAMI VE GEREKSİNİMLER

Amacımız, AVEVA PI ve GE Proficy gibi devlerle rekabet edebilecek, modern, mikroservis tabanlı bir **"Kurumsal Endüstriyel Historian Platformu"** inşa etmektir.

## DURUM ÖZETİ (2025-12-02)
*   **Aşama:** Planlama ve Çözümleme Tamamlandı. Uygulama (Phase 4) Başlıyor.
*   **Tamamlananlar:** PRD, Mimari, Epics & Stories.
*   **Sıradaki Adım:** Epic 0 (Altyapı) Kurulumu.

## TEMEL DİREKLER
1.  **Mimari:** **Polyglot Monorepo** yapısında, Docker üzerinde çalışan ve **NATS JetStream** ile haberleşen dağıtık mikroservisler.
2.  **Performans:** Kritik G/Ç (Ingestion/Engine) için **Rust**, İş Mantığı (Auth/Audit) için **Go**, UI için **React (Vite)**.
3.  **Uyumluluk:**
    *   **Alarm:** ISA 18.2 (Durum Makinesi).
    *   **FDA:** 21 CFR Part 11 (Re-authentication, Audit Trail).

## ANA MODÜLLER (KONTEYNERLER)
1.  **Ingestor (Rust):** Modbus/OPC-UA veri toplayıcı. Stateless. Protobuf ile NATS'a basar.
2.  **Engine (Rust):** LSM-Tree tabanlı TSDB. Gorilla Sıkıştırma. S3'e arşivleme.
3.  **Auth (Go):** Kimlik doğrulama servisi. JWT ve FDA Re-auth mantığı.
4.  **Audit (Go):** Değiştirilemez loglama (Chained Hash).
5.  **Sim (Python):** GEKKO tabanlı Dijital İkiz (Simülasyon).
6.  **Viz (React + Vite):** **uPlot** ile yüksek performanslı dashboard. State için **Zustand**.
7.  **Ops:** Docker Compose (Edge) ve Kubernetes (Cloud).

## TEKNİK KURALLAR VE DESENLER
*   **İletişim:**
    *   **Async:** NATS JetStream (Konu: `enterprise.site.area.line.device.sensor`).
    *   **Sync:** gRPC (Servisler arası), GraphQL (Frontend -> Backend).
    *   **Schema:** Tüm veri yapıları **Protobuf** ile tanımlanır (`crates/historian-core`).
*   **Frontend:**
    *   State Management: **Zustand** (Redux/Context kullanılmaz).
    *   Build: **Vite** (Next.js kullanılmaz).
*   **Hata Yönetimi:** Rust tarafında `thiserror` (Lib) ve `anyhow` (App). Asla `unwrap()` kullanılmaz.

## BAŞARI KRİTERLERİ
- **Yazma Hızı:** >50.000 olay/saniye (Tek düğüm).
- **Sorgu Hızı:** 1 yıllık veri <100ms.
- **Kaynak:** Edge cihazda <500MB RAM.
- **Güvenilirlik:** Ağ kesintisinde %0 veri kaybı (Store & Forward).
