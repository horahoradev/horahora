from selenium.webdriver.chrome.options import Options
from selenium import webdriver
from selenium.webdriver.common.by import By
import time

if __name__ == "__main__":
    chrome_options = webdriver.ChromeOptions()
    chrome_options.add_argument('--no-sandbox')
    chrome_options.add_argument('--window-size=1920,1080')
    chrome_options.add_argument('--headless')
    chrome_options.add_argument('--disable-gpu')
    driver = webdriver.Chrome(chrome_options=chrome_options)

    driver.get("http://nginx/authentication/login")

    time.sleep(5)
    assert driver.title == "Log in"

    time.sleep(5)

    username_box = driver.find_element(by=By.ID, value="auth-login-name")
    password_box = driver.find_element(by=By.ID, value="auth-login-password")
    submit_button = driver.find_element(by=By.ID, value="auth-login")
    username_box.send_keys("admin")
    password_box.send_keys("admin")
    submit_button.submit()
    time.sleep(5)
    print(driver.get_cookies())

    time.sleep(5)
    print(driver.current_url)
    print(driver.get_cookies())

    assert driver.current_url == "http://nginx/"
    assert driver.title == "Videos"

    driver.quit()
    print("wow we did it")
