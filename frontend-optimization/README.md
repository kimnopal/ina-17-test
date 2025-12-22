Untuk melakukan improvement dapat dilakukan beberapa hal berikut:

### Load Time

- Melakukan pagination di sisi server agar tidak meload data besar di awal
- Optimalisasi query (select data yang diperlukan saja, melakukan indexing)
- Melakukan caching untuk data yang jarang berubah

### Interactivity

- Implementasi debouncing pada search data
- Implementasi pagination
- Memoization

### Code Splitting

- Lazy loading pada halaman yang belum muncul di halaman pengguna
- Dynamic import untuk komponen yang berat
