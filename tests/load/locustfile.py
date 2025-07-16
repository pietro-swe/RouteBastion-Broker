from locust import HttpUser, task

class BrokerApiUser(HttpUser):
  @task
  def optimize(self):
    self.client.get("/v1/optimizations/sync")
