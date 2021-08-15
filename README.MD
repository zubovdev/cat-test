# Start

1. `docker-compose up`
2. `go build -o api`
3. `./api`

# Personal task management app API

I tried to design this API followed to the minimal requirements:

- ✗ - The requirement haven't been met
- ✓ - The requirement have been met
- ? - The requirement may have been met

```
Create user ✓
    * Minimal set of properties:
        * email ✓
        * first name ✓
        * last name ✓
Delete user ✓
Get user ✓
Get users ✓
Create a task ✓
    * task has to be assigned to a user ✓
    * task can occupy some certain amount of time ✓ (estimated_time)
    * task can’t overlap with other tasks ✗ (didn't get that)
    * reminder period ?✗ (didn't get that)
Delete task ✓
Get task ✓
Get tasks ✓
    * get list of the task ?✗ (didn't get that)
Reminder notification ✓?
    * by email ✓?
    * on-site ✗ (didn't get that)
```

You can also check swagger doc `docs/api.yml`.  
Postman collection: https://www.getpostman.com/collections/28efe41272c415841688  
Password for admin is: `Qwerty123`
