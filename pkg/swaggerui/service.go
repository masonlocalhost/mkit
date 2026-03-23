package swaggerui

import (
	"fmt"
	"io/fs"
	"net/http"
)

type Service struct {
	static fs.FS
}

func NewService() (*Service, error) {
	static, err := fs.Sub(SwaggerAsset, "asset")
	if err != nil {
		return nil, err
	}

	return &Service{
		static,
	}, nil
}

func (s *Service) GetStatic() http.FileSystem {
	return http.FS(s.static)
}

func (s *Service) GetInitializer(docUrl string) []byte {
	return []byte(fmt.Sprintf(`
		window.onload = function() {
		  //<editor-fold desc="Changeable Configuration Block">
		
		  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
		  window.ui = SwaggerUIBundle({
			url: "%s",
			dom_id: '#swagger-ui',
			deepLinking: true,
			presets: [
			  SwaggerUIBundle.presets.apis,
			  SwaggerUIStandalonePreset
			],
			plugins: [
			  SwaggerUIBundle.plugins.DownloadUrl
			],
			layout: "StandaloneLayout",
  			persistAuthorization: true,   // Keep auth info after reload
		    withCredentials: true,        // Allow sending cookies/credentials
		  });
		
		  //</editor-fold>
		};
	`, docUrl))
}
