package config

type Config struct {
	Logger Logger `json:"logger"`
	Server Server `json:"server"`
	Cache  Cache  `json:"cache"`
}

type Logger struct {
	Level            string   `json:"level"`
	OutputPaths      []string `json:"outputPaths"`
	ErrorOutputPaths []string `json:"errorOutputPaths"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Cache struct {
	Capacity string `json:"capacity"`
}
