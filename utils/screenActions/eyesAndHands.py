import pyautogui

# TODO: Just do the screen actions with pyautogui - finish this
# TODO: install opencv
location = pyautogui.locateOnScreen()
if location is not None:
    print("Found the ball!")
    # You can use pyautogui to click on the ball or perform other actions
    x, y = pyautogui.center(location)
    pyautogui.click(x, y)