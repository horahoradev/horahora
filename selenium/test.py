from selenium.webdriver.chrome.options import Options
from selenium import webdriver
from selenium.webdriver.common.by import By
import time

if __name__ == "__main__":
    chrome_options = Options()
    chrome_options.add_argument('--no-sandbox')
    chrome_options.add_argument('--window-size=1920,1080')
    chrome_options.add_argument('--headless')
    chrome_options.add_argument('--disable-gpu')

    driver = webdriver.Chrome(options=chrome_options)

    driver.get("http://nginx/authentication/login")

    time.sleep(1)
    assert driver.title == "Log in"

    time.sleep(1)

    username_box = driver.find_element(by=By.ID, value="auth-login-name")
    password_box = driver.find_element(by=By.ID, value="auth-login-password")
    submit_button = driver.find_element(by=By.ID, value="auth-login")
    username_box.send_keys("admin")
    password_box.send_keys("admin")
    submit_button.submit()
    time.sleep(1)

    assert driver.current_url == "http://nginx/"
    assert driver.title == "Videos"

    driver.get('http://nginx/api/home')
    assert len(driver.get_cookies()) >= 1

    driver.get('http://nginx/account/archives')
    time.sleep(1)
    video_box = driver.find_element(by=By.ID, value="new-video-url")
    video_form = driver.find_element(by=By.ID, value="new-video")
    video_box.send_keys("https://www.youtube.com/watch?v=GWAtnzcbfFQ")
    video_form.submit()

    time.sleep(1)
    dl_card = driver.find_element(by=By.ID, value="download-request-1")
    assert dl_card is not None

    # Search for author
    for i in range(1, 120):
        try:
            driver.get('http://nginx/search?search=%E3%82%82%E3%81%A1%E3%81%86%E3%81%A4%E3%81%AD&category=upload_date&order=desc')
            time.sleep(1)
            video_href = driver.find_element(by=By.ID, value="/videos/1")
            video_href.click()
            break
        except:
            print("Failed to find video 1 on page, attempt " + i)

    time.sleep(1)
    assert driver.current_url == "http://nginx/videos/1"
    assert driver.title == 'Video "おくすり飲んで寝よう / 初音ミク - もちうつね" (1) by "もちうつね" (1)'

    driver.quit()
    print("wow we did it")
