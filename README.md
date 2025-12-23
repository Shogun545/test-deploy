#  Team 11 ระบบบริหารกระบวนการนัดหมายและติดตามการให้คำปรึกษานักศึกษา สำหรับสำนักวิชาวิศวกรรมศาสตร์ 

## รายชื่อสมาชิกทีม

| รหัสนักศึกษา| ชื่อ - นามสกุล | ระบบย่อย |
|---------------|------------------|------------------|
| B6608019 | นางสาวเนตรนภัทร ชำนินอก | ระบบจัดการผู้ใช้และโปรไฟล์ & ระบบจองนัดหมาย  |
| B6628239 | นายกิตติธัช แช่มขุนทด | ระบบแจ้งปัญหา & ระบบจัดการ FAQ ของอาจารย์ที่ปรึกษา  |
| B6631536 | นายภูวิศ แสนตา  | ระบบบันทึกการให้คำปรึกษา &ระบบรายงานความคืบหน้าหลังให้คำปรึกษา  |
| B6631659 | นายวงศกร ยอดกลาง | ระบบตั้งค่าระบบและปฏิทินวิชาการ & ระบบจัดการเวลาว่างของอาจารย์  |
| B6639334 | นางสาวพิมพ์นารา อดุลจันทรศร | ระบบอนุมัติ/เสนอเวลาใหม่การนัดหมาย & ระบบแจ้งเตือน  |


#### Build Image ใหม่ หลังจากแก้โค้ด Build + Run ด้วย docker-compose

```bash
# build และ run ใหม่ทั้งหมด
docker-compose up --build -d
```

#### ตรวจสอบ container / Stop / Remove Container เก่า (ถ้า run อยู่แล้ว)

```bash
docker ps                       # ดู container ที่กำลังรัน
docker-compose logs -f          # ดู log ของทุก container
docker-compose logs -f backend  # ดูเฉพาะ backend
docker stop <container_id>      # หยุด container
docker rm <container_id>        # ลบ container

```

#### ตรวจสอบ container

```bash
docker-compose stop             # หยุดทุก container
docker-compose stop backend     # หยุดเฉพาะ backend
docker-compose stop frontend    # หยุดเฉพาะ frontend

docker-compose restart          # restart ทุก container
docker-compose restart backend  # restart เฉพาะ backend
docker-compose restart frontend # restart เฉพาะ frontend
```
#### Reset / Remove container

```bash
docker-compose down
#ลบ container และ network ที่สร้างโดย docker-compose แต่ ไม่ลบ image

```

#### Push image ไป Docker Hub

```bash
docker-compose build        # rebuild image
docker tag project-backend:latest netnaphat/project-backend:latest
docker tag project-frontend:latest netnaphat/project-frontend:latest
docker push netnaphat/project-backend:latest
docker push netnaphat/project-frontend:latest

```

