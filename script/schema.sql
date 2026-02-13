-- Perintah untuk membuat tabel 'users'
-- File ini bisa kamu letakkan di folder 'scripts/schema.sql' atau 'migrations/init.sql'

CREATE TABLE IF NOT EXISTS users (
    -- ID unik untuk setiap pengguna, otomatis bertambah (Auto Increment)
    id INT AUTO_INCREMENT PRIMARY KEY,
    
    -- Nama lengkap pengguna, tidak boleh kosong
    full_name VARCHAR(255) NOT NULL,
    
    -- Nama pengguna unik, tidak boleh ada yang sama di database
    username VARCHAR(100) NOT NULL UNIQUE,
    
    -- Email unik sebagai identitas utama login, tidak boleh ada yang sama
    email VARCHAR(255) NOT NULL UNIQUE,
    
    -- Tempat menyimpan password yang sudah di-hash (bcrypt)
    password VARCHAR(255) NOT NULL,
    
    -- Kolom otomatis yang mencatat kapan user mendaftar
    -- Kita atur agar database otomatis mengisi waktu saat ini (CURRENT_TIMESTAMP)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Kolom otomatis yang mencatat kapan data terakhir kali diubah
    -- Database akan otomatis memperbarui kolom ini setiap kali ada perubahan data (ON UPDATE)
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Kolom untuk Soft Delete (penghapusan aman)
    -- Jika diisi waktu, berarti user dianggap "terhapus" namun datanya masih ada di database
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Menambahkan Index pada email agar proses pencarian (saat login) lebih cepat
    INDEX idx_user_email (email)
);