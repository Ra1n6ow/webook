wrk.method="GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.32.3"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI3NDY5ODgsIlVpZCI6MTAsIlVzZXJBZ2VudCI6IkFwaWZveC8xLjAuMCAoaHR0cHM6Ly9hcGlmb3guY29tKSJ9.OVVF_q8qYWnj1GTG4c2TwAI6YDqIah4YOUAgFfcs2l4"