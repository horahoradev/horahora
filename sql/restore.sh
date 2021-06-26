psql --dbname=postgresql://admin:password@localhost:5432/scheduler < backup_scheduler*.sql
psql --dbname=postgresql://admin:password@localhost:5432/videoservice < backup_videoservice*.sql
psql --dbname=postgresql://admin:password@localhost:5432/userservice < backup_userservice*.sql
