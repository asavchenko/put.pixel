#define GLEW_STATIC
#include <GL/glew.h>
#include <GLFW/glfw3.h>
#include <stdio.h>
#include <string.h>
#include "ogl.h"
#include <iostream>
#include <chrono>
#include <thread>
using namespace std::chrono_literals;
namespace ogl
{
    static const int HEIGHT = 480;
    static const int WIDTH = 640;
    static const int XM    = WIDTH - 0;
	static const int YM     = HEIGHT - 0;

    GLFWwindow* window;
    unsigned char* pixelArr;
    GLuint buffers[1];
    int index;
    int lastX;
    int lastY;
    int lastWidth;
    int lastHeight;
   
    bool isExit() {
        int ret = glfwWindowShouldClose(window);

        return ret > 0;
    } 

    void putPixel(int x, int y, unsigned char color)
    {
        int idx = (x + y*WIDTH) * 3;
        if (idx < 0) {
            return;
        }
        if (idx+2 > HEIGHT*WIDTH*3-1) {
            return;
        }
        pixelArr[idx] = color;
        pixelArr[idx+1] = color;
        pixelArr[idx+2] = color;
    }

    void putPixelRGB(int x, int y, unsigned char r, unsigned char g, unsigned char b)
    {
        int idx = (x + y*WIDTH) * 3;
        if (idx < 0) {
            return;
        }
        if (idx+2 > HEIGHT*WIDTH*3-1) {
            return;
        }
        pixelArr[idx] = r;
        pixelArr[idx+1] = g;
        pixelArr[idx+2] = b;
    }

    static void key_callback(GLFWwindow* window, int key, int scancode, int action, int mods)
    {
        if (key == GLFW_KEY_ESCAPE && action == GLFW_PRESS) {
            glfwSetWindowShouldClose(window, GLFW_TRUE);
        }
    }

    void error_callback(int error, const char* description)
    {
        fprintf(stderr, "Error: %s\n", description);
    }

    void window_close_callback(GLFWwindow* window)
    {
        glDeleteBuffers(1, buffers);
        if (window != nullptr) {
            glfwDestroyWindow(window);
        }
    }

    bool init(bool isFullScreen)
    {
        if (!glfwInit())
        {
           return false;
        }
        glfwSetErrorCallback(error_callback);
        glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 2);
        glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 1);
        glfwWindowHint(GLFW_RESIZABLE, false);
        glfwWindowHint(GLFW_REFRESH_RATE, 60);
        // glfwWindowHint(GLFW_DOUBLEBUFFER, GL_FALSE);
        // glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);
        if (isFullScreen) {
            window = glfwCreateWindow(WIDTH, HEIGHT, "My Title", glfwGetPrimaryMonitor(), NULL);
        } else {
            window = glfwCreateWindow(WIDTH, HEIGHT, "My Title", NULL, NULL);
        }
        if (!window)
        {
            glfwTerminate();
            return false;
        }
        glfwMakeContextCurrent(window);
        glewExperimental = GL_TRUE;
        GLenum err = glewInit();
        if (GLEW_OK != err)
        {
            std::cerr << "Error: " << glewGetErrorString(err) << std::endl;
            glfwTerminate();
            return false;
        }

	    glfwSwapInterval(1);

        glGenBuffers(1, buffers);
        glBindBuffer(GL_PIXEL_UNPACK_BUFFER, buffers[0]);
        glBufferData(GL_PIXEL_UNPACK_BUFFER, WIDTH*HEIGHT*4, NULL, GL_DYNAMIC_DRAW);

        // glBindBuffer(GL_PIXEL_UNPACK_BUFFER, buffers[1]);
        // glBufferData(GL_PIXEL_UNPACK_BUFFER, WIDTH*HEIGHT*4, NULL, GL_DYNAMIC_DRAW);
        glfwGetWindowPos(window, &lastX, &lastY);
        glfwGetWindowSize(window, &lastWidth, &lastHeight);
        glfwSetKeyCallback(window, key_callback);
        glfwSetWindowCloseCallback(window, window_close_callback);
        GLubyte* pboPtr = (GLubyte*)glMapBuffer(GL_PIXEL_UNPACK_BUFFER, GL_WRITE_ONLY);
        if (pboPtr == NULL) { std::cout<< "null ptr" << std::endl;
            glfwTerminate();
            return false;
        }
        if (!glUnmapBuffer(GL_PIXEL_UNPACK_BUFFER)) {
            std::cout<< "unable to unmap" << std::endl;
            glfwTerminate();
            return false;
        }
        pixelArr = (unsigned char*)pboPtr;
        clearScreen();

        return true;
    }

    void draw(void (*cb)()) 
    {
        while (!glfwWindowShouldClose(window)) // Keep running
        {
            GLubyte* pboPtr = (GLubyte*)glMapBuffer(GL_PIXEL_UNPACK_BUFFER, GL_WRITE_ONLY);
            if (pboPtr == NULL) {
                std::cout<< "null ptr" << std::endl;
                return;
            }
            if (!glUnmapBuffer(GL_PIXEL_UNPACK_BUFFER)) {
                std::cout<< "unable to unmap" << std::endl;
                return;
            }
            pixelArr = (unsigned char*)pboPtr;
            cb();
            glDrawPixels(WIDTH, HEIGHT, GL_RGB, GL_UNSIGNED_BYTE, NULL);
            glfwSwapBuffers(window);
            glfwPollEvents();
        }
    }

    void clearScreen() 
    {
	    memset(pixelArr, 0, WIDTH*HEIGHT*3);
    }

    int getWindowHeight()
    {
        return HEIGHT;
    }

    int getWindowWidth()
    {
        return WIDTH;
    }

    void swapBuffers()
    {
        // glfwSwapBuffers(window);
        glFinish();
    }

    void close()
    {
        glfwTerminate();
    }
}