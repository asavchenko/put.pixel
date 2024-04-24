#include "ogl/ogl.h"
#include <random>
#include <string.h>
#include "characters/characters.h"
#include "characters/character.h"
#include <iostream>

void randomPixels();
void testMovingWord();
char text[] = "IT WORKS!";
int w = ogl::getWindowWidth();
int h = ogl::getWindowHeight();
unsigned char color = 200;
int fontSize = 64;
int textWidth = strlen(text) * (characters::getCharacterWidth(fontSize) + characters::getSpaceSizeBtwCharacters(fontSize));
int textHeight = characters::getCharacterHeight(fontSize) + characters::getLineSpaceSize(fontSize);
int y = (h + 2*textHeight) / 2;
int x = (w - textWidth) / 2;
int l = strlen(text);
Chr** chrs;
///////////////////////////////////////////////////////////////////////////////
int main(int argc, char **argv)
{
    std::cout<<"starting" << std::endl;
    std::cout<<"starting ogl" << std::endl;
    ogl::init(false);
    chrs = new Chr*[l];
    std::cout<<"creating text" << std::endl;
	for (int i = 0; i < l; i++) {
        char32_t r = text[i];
        std::cout<<"creating ch"<< i << std::endl;
		chrs[i] = new Chr(r, x, y, color);
        std::cout<<"done creating ch"<< i << std::endl;
        chrs[i]->setCharacterSize(fontSize);
		x += characters::getCharacterWidth(fontSize) + characters::getSpaceSizeBtwCharacters(fontSize);
	}
    std::cout<<"done creating text" << std::endl;
    ogl::draw(testMovingWord);
    for (int i = 0; i < l; i++) {
        delete chrs[i];
    }
    delete chrs;
    ogl::close();
    
    return 0;
}

void testMovingWord()
{
    for (int i = 0; i < l; i++)
    {
        chrs[i]->move(0, -1);
        if (chrs[i]->getY() < 0) {
            chrs[i]->setY(ogl::getWindowHeight() + chrs[i]->getHeight());
        }
    }
}

void randomPixels()
{
    for (int y = 0; y < ogl::getWindowHeight(); y++) 
    {
		for (int x = 0; x < ogl::getWindowWidth(); x++) 
        {
			ogl::putPixel(x, y, (unsigned char)(rand() & 0x00fe));
		}
	}
}

