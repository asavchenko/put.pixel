#include "character.h"

int Chr::getCharacterSize()
{
    return m_size;
}

Chr Chr::setCharacterSize(int size)
{
    switch (size)
    {
    // case 14:
        // for (int i = 0; i < 17; i++)
        // {
        //     /* code */
        // }
        // utf8::getShape(m_ch);
    default:
        scale(size);
    }
    m_size = size;

    return *this;
}

int Chr::getCharacterWidth()
{
    switch (m_size)
    {
    case 14:
        return 17;
    default:
        if (m_size > 14)
        {
            return 17 + m_size - 14;
        }
        if (m_size < 0)
        {
            return 3;
        }
        return 17 - (14 - m_size);
    }
}

int Chr::getCharacterHeight()
{
    switch (m_size)
    {
    case 14:
        return 20;
    default:
        if (m_size > 14)
        {
            return 20 + m_size - 14;
        }
        if (m_size < 0)
        {
            return 5;
        }
        return 20 - (14 - m_size);
    }
}

int Chr::getSpaceSizeBtwCharacters()
{
    return getCharacterWidth() / 9;
}

int Chr::getLineSpaceSize()
{
    return getCharacterWidth() / 6;
}

Chr::Chr(char32_t ch, int x, int y, char color)
{
    m_ch = ch;
    m_wH = ogl::getWindowHeight();
    m_wW = ogl::getWindowWidth();
    bool(*original)[17] = utf8::getShape(m_ch);
    m_shape = new bool*[20];
    for (int j = 0; j < 20; j++)
    {
        m_shape[j] = new bool[17];
        for (int i = 0; i < 17; i++)
        {
            m_shape[j][i] = original[j][i];
        }
    }
    m_x = x;
    m_y = y;
    m_px = x;
    m_py = y;
    m_color = color;
    m_size = 14;
    m_width = getCharacterWidth();
    m_height = getCharacterHeight();
}

int Chr::getWidth()
{
    return m_width;
}

int Chr::getHeight()
{
    return m_height;
}

bool Chr::isVisible()
{
    if (m_x < 0 && m_x + m_width < 0)
    {
        return false;
    }

    if (m_y > 0 && m_y + m_height < 0)
    {
        return false;
    }

    if (m_x > m_wW && m_x + m_width > m_wW)
    {
        return false;
    }

    if (m_y > m_wH && m_y + m_height > m_wH)
    {
        return false;
    }

    return true;
}

void Chr::move(int dx, int dy)
{
    hide();
    m_px = m_x;
    m_py = m_y;
    m_x += dx;
    m_y += dy;

    show();
}

void Chr::hide()
{
    for (int i = getCharacterHeight() - 1; i >= 0; i--)
    {
        for (int j = getCharacterWidth() - 1; j >= 0; j--)
        {
            if (m_shape[i][j])
            {
                ogl::putPixel(m_px + j, m_py - i, 0);
            }
        }
    }
}

int Chr::abs(int i)
{
    if (i >= 0)
    {
        return i;
    }

    return -i;
}

void Chr::show()
{
    draw(m_shape, m_x, m_y);
}

void Chr::draw(bool **shape, int x, int y)
{
    for (int i = getCharacterHeight() - 1; i >= 0; i--)
    {
        for (int j = getCharacterWidth() - 1; j >= 0; j--)
        {
            if (shape[i][j])
            {
                ogl::putPixel(x + j, y - i, m_color);
            }
        }
    }
}

void Chr::scale(int size)
{
    bool(*original)[17] = utf8::getShape(m_ch);
    int ow = getCharacterWidth();
    int oh = getCharacterHeight();
    m_size = size;
    int nw = getCharacterWidth();
    int nh = getCharacterHeight();
    float kw = (float)nw / (float)ow;
    float kh = (float)nh / (float)oh;
    bool **resized = new bool*[nh];
    if (kw > 1 && kh > 1)
    {
        for (int j = 0; j < nh; j++)
        {
            resized[j] = new bool[nw]; // resized[j] = make([]int, nw)
            for (int i = 0; i < nw; i++)
            {
                int y = (int)std::ceil(j / kh);
                int x = (int)std::ceil(i / kw);
                if (x >= ow)
                {
                    x = ow - 1;
                }
                if (y >= oh)
                {
                    y = oh - 1;
                }

                resized[j][i] = original[y][x];
            }
        }
        // free m_shape and then assign new shape
        if (m_shape)
        {
            for (int i = 0; i < ow; i++)
            {
                if (m_shape[i])
                {
                    delete[] m_shape[i];
                }
            }
            delete[] m_shape;
        }
        m_shape = resized;

        return;
    }
    for (int i = 0; i < nh; i++)
    {
        resized[i] = new bool[nw]; // make([]int, nw)
    }
    for (int j = 0; j < oh; j++)
    {
        for (int i = 0; i < ow; i++)
        {
            if (!original[j][i])
            {
                continue;
            }
            int y = std::ceil(j * kh);
            int x = std::ceil(i * kw);
            if (y >= nh)
            {
                y = nh - 1;
            }
            if (x >= nw)
            {
                x = nw - 1;
            }

            resized[y][x] = original[j][i];
        }
    }
    m_shape = resized;
}
