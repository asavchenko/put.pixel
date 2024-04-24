#include "characters.h"
namespace characters{
    int getCharacterWidth(int size)
    {
        switch (size){
        case 14:
            return 17;
        default:
            if (size > 14) {
                return 17 + size - 14;
            }
            if (size < 0) {
                return 3;
            }

            return 17 - (14 - size);
        }
    }

    int getCharacterHeight(int size) 
    {
        switch (size) {
        case 14:
            return 20;
        default:
            if (size > 14) {
                return 20 + size - 14;
            }
            if (size < 0) {
                return 5;
            }
            return 20 - (14 - size);
        }
    }

    int getSpaceSizeBtwCharacters(int size)
    {
        return getCharacterWidth(size) / 9;
    }

    int getLineSpaceSize(int size) 
    {
        return getCharacterWidth(size) / 6;
    }
}