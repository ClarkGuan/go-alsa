package alsa

// #include <alsa/asoundlib.h>
//
// static char *no_const(const char *s) { return (char *)s; }
import "C"

// ConfigTopDir 返回 ALSA 的顶层配置目录。
// 默认情况下为 `/usr/share/alsa`；
// 如果存在环境变量 `$ALSA_CONFIG_DIR` 则取环境变量的值（前提是使用绝对路径）
func ConfigTopDir() string {
	return C.GoString(C.no_const(C.snd_config_topdir()))
}
