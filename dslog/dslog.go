// Package dslog is a wrapper for the slog library.  The handler converts the
// slog messages to the format expected by DeepCode.  The format is:
//
// {
//  "timestamp":1516134303,
//  "level":"ERROR",
//  "trace_id":"regbhg93h8g93u4fh9u34hfgerg45",
//  "message":{},
//  "env": {
//    "host":"10.54.123.123",
//    "arch:"arm64"
//  }
//  "context":{
//    "event":"analysis"
//    "owner_id":"deepcode-ai",
//    "repository_id":"asgard",
//    "analyzer":"Go"
//   }
// }
//
// To initialize the logger:
//
// 	dslog.Initialize(dslog.Option{
//		Writer: os.Stdout,
//		Level:  dslog.LevelDebug,
//	})
//
// To log a message:
// 	dslog.Debug("hello world", slog.String("name", "deepcode-ai"), slog.Int("age", 1))
//
// The additional fields are optional and will be added to the context field.

package dslog
