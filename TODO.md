# Ưu điểm
## Tách biệt rõ ràng:
- Phân tách rõ ràng giữa các tầng (transport, service, repository)
- Dễ dàng test từng thành phần độc lập

## Dễ bảo trì:
- Code được tổ chức theo từng tính năng (book, user, auth...)
- Dễ dàng thêm mới API hoặc thay đổi logic nghiệp vụ

## Tích hợp nhiều công nghệ:
- Sử dụng GORM cho ORM
- Redis cho caching
- MinIO cho lưu trữ file
- JWT cho xác thực
- Casbin cho phân quyền

## Xử lý lỗi tập trung:
- Có middleware xử lý lỗi chung
- Logging chi tiết

# Nhược điểm
## Độ phức tạp:
- Có thể quá phức tạp với các dự án nhỏ
- Nhiều lớp trung gian có thể gây khó hiểu cho người mới

## Khởi tạo phức tạp:
- Nhiều dependency cần khởi tạo (DB, Redis, MinIO...)
- Cần xử lý lỗi cẩn thận khi khởi tạo

## Performance:
- Nhiều lớp trừu tượng có thể ảnh hưởng hiệu năng
- Cần cân nhắc sử dụng connection pooling

## Documentation:
- Cần thêm tài liệu mô tả luồng xử lý
- Thiếu mô tả chi tiết về các API

# Đề xuất cải thiện
## Functionality
- [ ] pre-signed URL (rate-limit, check file type và size, validation về tên file, thêm bucket policy, đảm bảo file đúng content-type như đã khai báo)
- [ ] Thêm context withTimeout
- [ ] Khi globalrecover thì thêm thông báo tele hoặc email
## Performance:
- [x] Cân nhắc sử dụng connection pooling
- [ ] Redis Caching
- [ ] Gắn context vào trace-ID

## Tài liệu:
- [x] Thêm API documentation với Swagger
- [x] Viết README hướng dẫn cấu hình và chạy
- [ ] Thêm Param, Security vào swagger

## Testing:
- [ ] Thêm unit test và integration test
- [ ] coverage.html
- [ ] Sử dụng mock cho các dependency

## Cấu hình:
- [ ] Sử dụng environment variables thay vì file config cứng
- [ ] Thêm validation cho cấu hình
- [ ] Casbin with MySQL

## Giám sát:
- [ ] Thêm trace-ID for restapi, grpc
- [ ] Thêm metrics và monitoring
- [ ] Tích hợp với các công cụ theo dõi hiệu năng
- [ ] Thêm pprof

## Bảo mật:
- [ ] Thêm external, internal ID
- [ ] Thêm rate limiting
- [ ] Xử lý CORS đúng cách
- [ ] Validate input kỹ hơn
- [ ] Thêm authorities vào token, frontend sẽ dựa vào đó để tạo các nút, component theo authorities
- [ ] Thêm các API để quản lý các authorities
- [ ] Thêm các API để quản lý các roles

