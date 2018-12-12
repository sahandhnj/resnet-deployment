from locust import HttpLocust, TaskSet, task


class UserTasks(TaskSet):

    @task
    def predict(self):
        with open('input.jpg', 'rb') as git simage:
            self.client.post(
                "/predict?token=kbaHGfnd0XeQSOk0OL1eFdOkLSHdhp44tGPPZGw0D4rAtlg0cwx1gUQ4oij",
                data={},
                files={'file': image}
            )

class WebsiteUser(HttpLocust):
    min_wait = 2000
    max_wait = 5000
    task_set = UserTasks