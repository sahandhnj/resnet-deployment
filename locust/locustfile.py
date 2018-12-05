from locust import HttpLocust, TaskSet, task


class UserTasks(TaskSet):

    @task
    def predict(self):
        with open('input.jpg', 'rb') as image:
            self.client.post(
                "/predict",
                data={},
                files={'file': image}
            )

class WebsiteUser(HttpLocust):
    min_wait = 2000
    max_wait = 5000
    task_set = UserTasks