from locust import FastHttpUser, task

class BrokerApiUser(FastHttpUser):
  @task
  def optimize(self):
    self.client.get("/v1/optimizations/sync")
