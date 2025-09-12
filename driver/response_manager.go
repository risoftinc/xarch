package driver

import (
	"log"

	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/config"
	"go.risoftinc.com/xarch/constant"
)

func ResponseManager(cfg config.ResponseManager) (*goresponse.ResponseConfig, error) {
	responseConfig, err := goresponse.LoadConfig(goresponse.ConfigSource{
		Method: cfg.Method,
		Path:   cfg.Path,
	})

	if err != nil {
		return nil, err
	}

	customMessageTemplateBuilder(responseConfig)

	return responseConfig, nil
}

func ResponseManagerAsync(cfg config.ResponseManager) (*goresponse.AsyncConfigManager, error) {
	asyncManager := goresponse.NewAsyncConfigManager(goresponse.ConfigSource{
		Method: cfg.Method,
		Path:   cfg.Path,
	}, cfg.Interval)

	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		log.Printf("Config response manager updated!")
	})

	if err := asyncManager.Start(); err != nil {
		return nil, err
	}

	customMessageTemplateBuilder(asyncManager.GetConfig())

	return asyncManager, nil
}

func customMessageTemplateBuilder(responseConfig *goresponse.ResponseConfig) {
	responseConfig.AddMessageTemplates(
		// Connection errors -> 503/UNAVAILABLE
		goresponse.NewMessageTemplateBuilder(constant.ErrorConnectionRefused).
			WithTemplate("Database connection refused. The database server may be down or unreachable.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Database connection refused. The database server may be down or unreachable.",
				constant.IdLanguage: "Koneksi ke database ditolak. Server database mungkin sedang down atau tidak dapat diakses.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 503,
				constant.ProtocolGrpc:   14,
			}).
			Build(),

		goresponse.NewMessageTemplateBuilder(constant.ErrorTooManyConnections).
			WithTemplate("Too many connections to database. Please try again in a few moments.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Too many connections to database. Please try again in a few moments.",
				constant.IdLanguage: "Terlalu banyak koneksi ke database. Silakan coba lagi dalam beberapa saat.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 503,
				constant.ProtocolGrpc:   14,
			}).
			Build(),

		goresponse.NewMessageTemplateBuilder(constant.ErrorConnectionTimeout).
			WithTemplate("Database connection timeout. Please check your network connection or try again.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Database connection timeout. Please check your network connection or try again.",
				constant.IdLanguage: "Koneksi ke database mengalami timeout. Periksa koneksi jaringan atau coba lagi.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 503,
				constant.ProtocolGrpc:   14,
			}).
			Build(),

		goresponse.NewMessageTemplateBuilder(constant.ErrorDnsError).
			WithTemplate("Cannot resolve database host. Please check DNS configuration or server address.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Cannot resolve database host. Please check DNS configuration or server address.",
				constant.IdLanguage: "Tidak dapat menemukan host database. Periksa konfigurasi DNS atau alamat server.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 503,
				constant.ProtocolGrpc:   14,
			}).
			Build(),

		// Auth errors -> 401/UNAUTHENTICATED
		goresponse.NewMessageTemplateBuilder(constant.ErrorAuthFailed).
			WithTemplate("Authentication failed. Please check the database credentials being used.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Authentication failed. Please check the database credentials being used.",
				constant.IdLanguage: "Autentikasi gagal. Periksa kredensial database yang digunakan.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 401,
				constant.ProtocolGrpc:   16,
			}).
			Build(),

		goresponse.NewMessageTemplateBuilder(constant.ErrorAccessDenied).
			WithTemplate("Access denied. Please ensure the user has sufficient permissions to access the database.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Access denied. Please ensure the user has sufficient permissions to access the database.",
				constant.IdLanguage: "Akses ditolak. Pastikan pengguna memiliki izin yang cukup untuk mengakses database.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 401,
				constant.ProtocolGrpc:   16,
			}).
			Build(),

		// Internal errors -> 500/INTERNAL
		goresponse.NewMessageTemplateBuilder(constant.ErrorDriverError).
			WithTemplate("Database driver error occurred. Please check driver configuration or contact administrator.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "Database driver error occurred. Please check driver configuration or contact administrator.",
				constant.IdLanguage: "Terjadi kesalahan pada driver database. Periksa konfigurasi driver atau hubungi administrator.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 500,
				constant.ProtocolGrpc:   13,
			}).
			Build(),

		goresponse.NewMessageTemplateBuilder(constant.ErrorSslTlsError).
			WithTemplate("SSL/TLS connection error. Please check certificate or database security configuration.").
			WithTranslations(map[string]string{
				constant.EnLanguage: "SSL/TLS connection error. Please check certificate or database security configuration.",
				constant.IdLanguage: "Kesalahan pada koneksi SSL/TLS. Periksa sertifikat atau konfigurasi keamanan database.",
			}).
			WithCodeMappings(map[string]int{
				constant.ProtocolWebApi: 500,
				constant.ProtocolGrpc:   13,
			}).
			Build(),
	)
}
