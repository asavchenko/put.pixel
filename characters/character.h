#include "utf8.h"
#include <cmath>
#include "../ogl/ogl.h"
class Chr
{
    char32_t m_ch;
    int m_size;
    int m_x;
    int m_y;
    int m_px;
    int m_py;
    char m_color;
    bool **m_shape;
    int m_wH;
    int m_wW;
    int m_width;
    int m_height;

public:
    int getX() { return m_x;}
    int getY() { return m_y;}
    Chr setX(int x) { m_x = x; return *this;}
    Chr setY(int y) { m_y = y; return *this;}
    int getCharacterSize();
    Chr setCharacterSize(int size);
    int getCharacterWidth();
    int getCharacterHeight();
    int getSpaceSizeBtwCharacters();
    int getLineSpaceSize();
    Chr(char32_t ch, int x, int y, char color);
    int getWidth();
    int getHeight();
    bool isVisible();
    void move(int dx, int dy);
    void hide();
    int abs(int i);
    void show();
    void draw(bool **shape, int x, int y);
    void scale(int size);
};
