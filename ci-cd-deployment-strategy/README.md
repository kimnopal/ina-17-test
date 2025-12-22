# Deployment Strategy

### 1. Branch Strategy

Terdapat 2 branch yaitu main untuk production dan dev untuk staging.

### 2. Proses CI

Step:

- Install dependensi
- Melakukan testing
- Melakukan linting
- Build Image (multi stage build, tagging dengan hash commit)

Proses build image dilakukan terakhir setelah seluruh aplikasi sudah memenuhi standar, karena proses build image membutuhkan waktu dan resource. Sehingga ketika proses linting atau testing gagal, maka tidak akan membuang waktu dan resource.

### 3. Proses CD

Step:

- Push image ke regisitry
- CD pipeline akan pull image dari registry
- Deploy ke production (blue green deployment)

### 4. Post Deployment

Setelah berhasil melakukan deployment, tidak lupa untuk melakukan monitoring, healthcheck, logging, dan rollback jikalau ada masalah
