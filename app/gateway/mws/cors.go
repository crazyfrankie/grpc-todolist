package mws

//func CORS(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
//		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, x-jwt-token")
//		w.Header().Set("Access-Control-Allow-Credentials", "true")
//		w.Header().Set("Access-Control-Max-Age", "43200") // 12小时
//
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		// 继续执行下一个处理程序
//		next.ServeHTTP(w, r)
//	})
//}
