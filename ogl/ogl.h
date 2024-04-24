#include <GLFW/glfw3.h>

namespace ogl
{
    int getWindowHeight();
    int getWindowWidth();
    void putPixel(int x, int y, unsigned char color);
    void putPixelRGB(int x, int y, unsigned char r, unsigned char g, unsigned char b);
    bool init(bool isFullScreen);
    void clearScreen();
    void draw(void (*run)());
    bool isExit();
    void close();
}