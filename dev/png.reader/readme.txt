There are many articles on the Internet about PNG format. Some are good, some are not. This attempt is to provide
my own understanding of the format. I hope it will be useful to someone. With this being said let's jump into it.

High level overview
-------------------
I want to start with a picture showing the structure of a PNG file.

           -----------------------------------------------------------------
           |8 bytes HEADER
           -----------------------------------------------------------------
           |                            ... CHUNKS ...
           -----------------------------------------------------------------
           | There numerous types of chunks. We are interested
           |   in IHDR, IDAT, IEND, PLTE
           | IHDR usually comes first, then PLTE, then IDAT, then IEND
           | While all the chunks are different,
           | they have in common this structure:
           |    4 bytes length
           |    4 bytes type (it literally contains these letters IDAT, IEND, etc)
           |    data, you know its length from the first field
           |    4 bytes CRC
           -----------------------------------------------------------------
           | IHDR chunk data is 13 bytes long
           | 4 bytes width
           | 4 bytes height
           | 1 byte bit depth
           | 1 byte color type
           | 1 byte compression method
           | 1 byte filter method
           | 1 byte interlace method
           | OK, as you can see, there is some useful information here!
           -----------------------------------------------------------------
           | PLTE chunk it's not always present. So let's skip it for now
           -----------------------------------------------------------------
           | IDAT chunk is the most interesting one. It contains the image data.
           | This is not a single chunk. It can be split into multiple chunks.
           | See below:
           -----------------------------------------------------------------
           | IDAT1 |4 bytes length| 4 bytes type| data | 4 bytes CRC
           -----------------------------------------------------------------
           | IDAT2 |4 bytes length| 4 bytes type| data | 4 bytes CRC
           -----------------------------------------------------------------
           | IDAT3 |4 bytes length| 4 bytes type| data | 4 bytes CRC
           -----------------------------------------------------------------
           .
           .
           .
           | IDATN |4 bytes length| 4 bytes type| data | 4 bytes CRC
           -----------------------------------------------------------------
           | So to combine the image data you need to concatenate the data fields
           | from all IDAT chunks
           |
           | The process is next:
           | 1. Read the length of the chunk
           | 2. Read the type of the chunk to make sure that it's IDAT
           | 3. Read the data, add it to the image data
           | imgData = IDAT1.data + IDAT2.data + ... + IDATN.data
           | So after you combine all IDAT chunks you will have the image data
           -----------------------------------------------------------------
           |
           .
           .                  Optional chunks
           .
           |
           -----------------------------------------------------------------
           | IEND chunk is the last one. It's 0 bytes long
           -----------------------------------------------------------------

Doesn't look scary? Yeah, it's not that bad. Let's go through the details.
Imagine that we have read the image data into an array of bytes:

    imgData = IDAT1.data + IDAT2.data + ... + IDATN.data

Ok, so far so good. The hardest part is in that this data is compressed using DEFLATE. This is the method used
not only by PNG, but also by gzip. We will go into details on how to uncompress the imgData, but for now let's
pretend that this is already done:

    uncompressedImgData = zlib.Decode(imgData)

If you think that you are done and can move on to the outputting it to the screen, you are wrong!
The uncompressedImgData has next structure (I will use RGBA color type):

    ----------------->----------------------------------------------->------------------------------
    |Filter byte|Red 1 byte|Green 1 byte|Blue 1 byte|Alpha 1 byte|... repeat <image width> times   |
    |only one   |          |            |           |            |                                 |
    |filter byte|          |            |           |            |                                 V
    |per line   |          |            |           |            |                                 |
    -------------------------------<---------------------------------------------<------------------
    |
    V----------------------------------------------->------------------------------------------->
    |Filter byte|Red|Green|Blue|Alpha|R|G|B|A|... repeat <image width> times
    --------------------------------------------------------------------------------------------
    .
    V                                      <image height> times
    .
    --------------------------------------------------------------------------------------------
    |Filter byte|Red|Green|Blue|Alpha|
    --------------------------------------------------------------------------------------------


Notice the Filter byte? It's a byte that is added before each scan line. It's used to improve the compression.
But we are very close to outputting the image to the screen. See the next picture (it shows the direction
of how to go through the image data):

    -->------
            |
            V
            |
    ---<-----
    |
    V
    |
    -----------------------
    | | | | | | | | | | |  |
    ------------------------
    | | | | | | | | | | |  |
    ------------------------
    | | | | | | | | | | |  |
    ------------------------
    .
    .
    .
    ------------------------
    | | | | | | | | | | |  |
    ------------------------

You start from the top left or in terms of uncompressedImgData it will be uncomressedImgData[0]. Then you go to
the right until you reach the end of the line. Then you go to the next line and so on. Or using pseudo code:

    for y = 0 to height
        for x = 0 to width
            r = uncompressedImgData[y * width * 4 + x]
            g = uncompressedImgData[y * width * 4 + x + 1]
            b = uncompressedImgData[y * width * 4 + x + 2]
            a = uncompressedImgData[y * width * 4 + x + 3]
            output pixel to the screen:
                PutPixel(x, y, r, g, b, a)

Why there is 4? Because for each pixel we have 4 bytes: r, g, b and a.

And note that I ignored the filter byte. Pretended that the uncompressedImgData doesn't have it. It means that
we are really close to outputting the image to the screen. But we are not there yet. Ok let's deal with the
filter byte. It is added before each scan line.

Possible values for the filter byte are:

    0 - None
    1 - Sub
    2 - Up
    3 - Average
    4 - Paeth

if we see None  it means that the line that follows can be copied without any modifications, using pseudo code:
    for y = 0 to height
        idx = y * width * 4 + y
        filterByte = uncompressedImgData[idx]               // filter byte is the first byte in the line
                                                            // 1 + width * 4
                                                            // 1 + width * 4
                                                            // ...
                                                            // 1 + width * 4

            if filterByte  == 0
                for x = 0 to width * 4
                    filteredImgData[y * width * 4 + x] = uncompressedImgData[idx + x] // just copy each byte
                                                                                      // from the
                                                                                      // uncompressedImgData
if we see Sub:
            if filterByte == 1
                for x = 0 to width * 4
                    if x < 4
                        prev = 0
                    else
                        prev = filteredImgData[y * width * 4 + x - 4]
                    filteredImgData[y * width * 4 + x] = uncompressedImgData[idx + x] + prev
if we see Up:
            if filterByte == 2
                for x = 0 to width * 4
                    if y == 0
                        filteredImgData[y * width * 4 + x] = uncompressedImgData[idx + x]
                    else
                        upperByte = filteredImgData[(y - 1) * width * 4 + x]
                        filteredImgData[y * width * 4 + x] = uncompressedImgData[idx + x] + upperByte
if we see Average:
            if filterByte == 3
                for x = 0 to width * 4
                    prev = 0
                    apper = 0
                    average = 0
                    if y != 0
                        apper = filteredImgData[(y - 1) * width * 4 + x]
                    if x > 4
                        prev = filteredImgData[y * width * 4 + x - 4]
                    average = (prev + upperByte) >> 1 //  divide by 2

                    filteredImgData[y * width * 4 + x] = uncompressedImgData[idx + x] +  average
if we see Paeth:
            if filterByte == 4
                for x = 0 to width * 4
                    a = 0
                    b = 0
                    c = 0
                    if y != 0
                        a = filteredImgData[(y - 1) * width * 4 + x]
                        if x > 4
                            c = filteredImgData[(y - 1) * width * 4 + x - 4]
                    if x > 4
                        b = filteredImgData[y * width * 4 + x - 4]

                    filteredImgData[y * width * 4 + x] = uncompressedImgData[idx + x] + PaethPredictor(a, b, c)

            where PaethPredictor is:
                PaethPredictor(a, b, c)
                    p = a + b - c
                    pa = abs(p - a)
                    pb = abs(p - b)
                    pc = abs(p - c)
                    if pa <= pb and pa <= pc
                        return a
                    else if pb <= pc
                        return b
                    else
                        return c

So what we are doing here is we need to apply filters to the bytes of the data. As you can see what it means is
applying certain formulas to the bytes using already filtered data. Not as hard as it looks:
    |FILTER|BYTES.....
    F = SOME FUNCTION (based on the FILTER) we need to apply to each byte from the BYTES on the line
    so it becomes:
    F(BYTE1)|F(BYTE2)|F(BYTE3)|F(BYTE4)|F(BYTE5)|F(BYTE6)|F(BYTE7)|F(BYTE8)|F(BYTE9)|F(BYTE10)|...
There is one important detail though:
    if we add two bytes we can go beyond byte limits. So for example everything inside the PaethPredictor
    needs to use 16-bit integers to avoid overflow.

And not only the PaethPredictor but in general the result needs to fit into byte boundaries,
i.e. use modulus 256 arithmetic:
    when we have a negative number we need to add 256 to it:
        RESULT = (BYTE1 + BYTE2) % 256
    or in general:
        RESULT = F(BYTE) % 256

Ok, so after we have applied the filters we are ready to output the image to the screen. Using the earlier
pseudo code:
        for y = 0 to height
            for x = 0 to width
                r = filteredImgData[y * width * 4 + x]
                g = filteredImgData[y * width * 4 + x + 1]
                b = filteredImgData[y * width * 4 + x + 2]
                a = filteredImgData[y * width * 4 + x + 3]
                output pixel to the screen:
                    PutPixel(x, y, r, g, b, a)


Now this all should be correct - the filter byte is no longer ignored here.
================================================================================================
What about the interlaced images?
We will not have a single rectangle of bytes in this case. An interlaced image is a series of
reduced images. We will have only portions of the whole image. There are 7 passes to go trough,
7 reduced images. If the whole image has 100% of data, a reduced image contains only N % of the whole
data, something like:
    allPixels = pixels from pass 1 + pixels from pass2 ... + pixels from pass7

How the pixels are selected to form a reduced image? Here is the pattern that applies to the entire image:

    1 6 4 6 2 6 4 6
    7 7 7 7 7 7 7 7
    5 6 5 6 5 6 5 6
    7 7 7 7 7 7 7 7
    3 6 4 6 3 6 4 6
    7 7 7 7 7 7 7 7
    5 6 5 6 5 6 5 6
    7 7 7 7 7 7 7 7

so on the first pass only these pixels are selected:

    (0, 0),   (8, 0),   (16, 0), ...,   (N*8, 0)
    (0, 8),   (8, 8),   (16, 8), ...,   (N*8, 8)
    ...
    (0, M*8), (8, M*8), (16, M*8), ..., (N*8, M*8)

    where N = (imgWidth  - 1) / 8 + 1
          M = (imgHeight - 1) / 8 + 1

These reduced images have own FILTERS:
        FILTER BYTE| PIXEL DATA | PIXEL DATA...
        FILTER BYTE| PIXEL DATA | PIXEL DATA...
        ...
        FILTER BYTE| PIXEL DATA | PIXEL DATA...

so we apply the same FILTER logic described earlier to get the RAW PIXEL DATA (or original pixel data)
Then we insert these filtered PIXELS back into their corresponding cells in the original image pixel array:

    after the first pass we should have only these pixels:
    --------------------------------------------------
    |0|          |0|          ...           |0|       |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |0|          |0|          ...           |0|       |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |
    --------------------------------------------------
    |                                                 |

and so on, with every new pass we get more and more pixels, until the entire image is reconstructed

================================================================================================
Now let's go back to the DEFLATE algorithm. Yes, this step can be ignored if you use a library that
decompresses the data for you. But if you want to understand how it works, let's dive into it.

I want to start with the structure again:

        As you remember we stopped at the point where we have the imgData as follows:
            imgData = IDAT1.data + IDAT2.data + ... + IDATN.data

        What imgData consists of is:
            imgData:
                CM (CompressionMethod - 4 bits)
                -------------------------------
                CINFO (Compression info - 4 bits)
                -------------------------------
                FLG (FLaGs - 8 bits)
                    FCHECK (5 bits)
                    FDICT (1 bit)
                    FLEVEL (2 bits)
                -------------------------------
                if FDICT is set then we have a dictionary
                    4 bytes Adler32
                if FDICT is not set then these 4 bytes are not present
                -------------------------------
                BLOCK1 + BLOCK2 + ... + BLOCKN
                -------------------------------
                CHECKSUM is 2 bytes long
                -------------------------------

                where BLOCK has next structure:
                    HEADER + DATA
                        HEADER is 3 bytes long:
                            BFINAL (1 bit)
                            BTYPE (2 bits)
                        DATA is variable length (more about it later)
                if BFINAL is set then it's the last block (we reached BLOCKN from the illustration above)

CM  is always = 8
CINFO is always = 7

So you can check your code at this point if you have these values equal to the 8 and 7 accordingly.
If not, you are doing something wrong.

FLG is a byte that consists of 8 bits. The first 5 bits are FCHECK, the next bit is FDICT and the last 2 bits
are FLEVEL.

FCHECK is a 5-bit field that is used to check the integrity of the data. It's calculated as follows:
    FCHECK = CM * 256 + CINFO

FDICT when set to 1 means that we have a DICT dictionary identifier immediately after the FLG byte. Let's not worry
about it for now. Still it needs to be checked to determine if we need to skip the 4 bytes right after the FLG byte.

FLEVEL is a 2-bit field that is used to indicate the compression level. Safe to ignore it.

-----------------------
The fun part begins here. The BLOCK is a sequence of bits (not bytes!). Just to remind you that each
    BLOCK has a HEADER and the DATA portions.
HEADER is a 3 bit long field. The first bit is BFINAL and the next 2 bits are BTYPE.
    if BFINAL is set to 1 it's the last block.
    BTYPE is used to indicate the type of the block.
    There are 3 types of blocks:
        0 - No compression
        1 - Fixed Huffman codes
        2 - Dynamic Huffman codes
        3 - Reserved (if you read BTYPE = 3 then you are doing something wrong)

----------------------------
No compression block
----------------------------
The easiest one. We simply need to read the data and copy it to the output. But because we are dealing with bits,
not bytes, we need to be careful. The data is byte-aligned. It means that the block starts with a byte boundary.
So if we have some bits left from the previous block, we need to skip them and start reading the data from the
next byte.

    Illustration:
               we have read ...3 bits of the BLOCK HEADER but there are possible:
                         ---+--------+--------+--------+--------+==================================+
... some remaining bits     | LEN    |  LEN   | NLEN   |  NLEN  |  ... LEN bytes of literal data...|
                         ---+--------+--------+--------+--------+==================================+
       and we read data bit by bit. If imagine data as as a stream of bits, we consume this stream
       one bit at a time. The direction of the stream:

       <---------------------------------<----------------------<-------------------


Ok, if we work on a bit level why do we need to account for the byte boundaries? Let me illustrate:

        Let's say we have the following stream of data:

            +--------+--------+--------+--------+--------+--------+
            | BYTE   |  BYTE  | BYTE   | BYTE   |  BYTE  | BYTE   |
            +--------+--------+--------+--------+--------+--------+

<---------------------------------<----------------------<------------------- STREAM DIRECTION
when we consume it we move from left to right through the data bit by bit, so our position in the stream
is not aligned with the byte boundaries:

                        |
            +--------+--------+--------+--------+--------+--------+
            | BYTE   |  BYTE  | BYTE   | BYTE   |  BYTE  | BYTE   |
            +--------+--------+--------+--------+--------+--------+
                        |
                        Our position in the stream
So what can happen is
                        |
            +--------+--------+--------+--------+--------+--------+
            | BYTE   |  BYTE  | BYTE   | BYTE   |  BYTE  | BYTE   |
            +--------+--------+--------+--------+--------+--------+
                        |
                        Our position in the stream
                        ||||LEN 8 bits|LEN 8 bits| NLEN 8 bits| NLEN 8 bits|... LEN bytes of data...
                        3 bits HEADER

What documentation states that you actually can't do that and need to account for the byte boundaries:
                        |
            +--------+--------+--------+--------+--------+--------+
            | BYTE   |  BYTE  | BYTE   | BYTE   |  BYTE  | BYTE   |
            +--------+--------+--------+--------+--------+--------+
                        |
                        Our position in the stream
                        ||||skip to the next byte and start reading the data from the next byte
                        3 bits HEADER

                       read 3 bit of the block's HEADER, then skip to the next byte
                       |
            +--------+--------+--------+--------+--------+--------+---------------------------------------------
            | BYTE   |  BYTE  | BYTE   | BYTE   |  BYTE  | BYTE   |                        |||| ..next block data
            +--------+--------+--------+--------+--------+--------+---------------------------------------------
                       ||||...|  LEN   |  LEN   | NLEN   | NLEN   |... LEN bytes of data...||||
                       Our position in the stream                                          3 bits of the next
                                                                                           block HEADER

I hope it makes sense now.
NLEN is the one's complement of LEN. It's used to check the integrity of the data. If LEN and NLEN don't match
then there is an error in the data. But we can ignore this for now.

--------------------------------
Dynamic Huffman codes
--------------------------------
Let's skip Fixed Huffman codes as we need to understand Dynamic Huffman codes first. The Fixed Huffman
codes are a subset of the Dynamic Huffman codes. Once you understand the Dynamic Huffman codes, you will understand
the Fixed Huffman codes as well.

As usual let's start with the structure of the block:

    HEADER + DATA
WHERE DATA begins with the following three integers:
    HLIT (5 bits) - number of literal/length codes - 257
    HDIST (5 bits) - number of distance codes - 1
    HCLEN (4 bits) - number of code length codes - 4
    ...
A lot of new concepts here. It means it's time to talk how the data is compressed.
literal/length, distance concepts refer to LZ77 algorithm.

LZ77's idea is to replace repeated sequences of bytes with a reference to the previous occurrence of the same
sequence. The reference consists of two parts: the length of the sequence and the distance from the current
position to the previous occurrence of the sequence.

Sounds a bit confusing so let's illustrate it:

    Suppose we have the following data:
        ABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA

A lot's of letters A. We can compress them using LZ77 algorithm. The compressed data will look like this:
    AB<65, 1>
    where 65 is the length of the sequence and 1 is the distance to the previous occurrence of the sequence.
    A is the byte that is repeated 65 times.

AB651 is much shorter than the original data. Ok if the idea in general more or less is clear, there are some
problems with it. How do we know that 65 means the length of the sequence and 1 is the distance to the previous
occurrence? How to distinguish between the length of the sequence and the distance to the previous occurrence
and the actual data? Good question! We need a map, or it's called alphabet in the LZ77 algorithm.

Alphabet:
    the first 255 characters (symbols/codes) actually still map to the 255 literals of the ASCII table
	the 256 symbol/code means "stop token"
	257 - 285 indicate "length codes/duplication length"

And it means that we need to go beyond the byte boundaries. If you wandered why do we need to work with bits and
not bytes, here is the answer.

If the first 255 symbols in the alphabet represent the actual symbols. Letters A and B from our example,
it's time to talk about the rest of the symbols. The symbols from 257 to 285. 256 is a special symbol, we will
cover it later.

	     Extra                 Extra               Extra
	Code Bits Length(s)   Code Bits Lengths   Code Bits Length(s)
	---- ---- ------      ---- ---- -------    ---- ---- -------
	257   0      3        267   1   15,16     277   4   67-82
	258   0      4        268   1   17,18     278   4   83-98
	259   0      5        269   2   19-22     279   4   99-114
	260   0      6        270   2   23-26     280   4  115-130
	261   0      7        271   2   27-30     281   5  131-162
	262   0      8        272   2   31-34     282   5  163-194
	263   0      9        273   3   35-42     283   5  195-226
	264   0     10        274   3   43-50     284   5  227-257
	265   1   11,12       275   3   51-58     285   0    258
	266   1   13,14       276   3   59-66

I'll try to explain this table with an example. Let's say we have encountered the symbol 257. It means that
    the length of the data to be copied is 3.
258 means 4, 259 5 and so on, note it starts from the minimal length of 3. Otherwise there is no sense in
    using the LZ77algorithm. If the repeated sequence is shorter than 3 symbols it can't be compressed.

Ok now what if we have encountered 265, see the extra bits column, it says 1. It means that we need to read
1 extra bit and the actual length will be
    11 + 1 = 12 if this extra bit is set,
or
    11 + 0 = 11 if this extra bit is not set.

The same idea applies to the rest of lengths codes. We read extra bits, convert them to a number, then
add it to the base length.

Let's discuss the "stop token" - the 256 symbol. When it's reached it means the end of the BLOCK. So you need to
proceed with reading the HEADER of the next BLOCK.

Ok, what about the map for distances?

         Extra                 Extra                   Extra
    Code Bits Dist        Code Bits Dist          Code Bits Dist
    ---- ---- ----        ---- ---- ----          ---- ---- ----
    0    0    1           10   4    33-48         20   9    1025-1536
    1    0    2           11   4    49-64         21   9    1537-2048
    2    0    3           12   5    65-96         22   10   2049-3072
    3    0    4           13   5    97-128        23   10   3073-4096
    4    1    5,6         14   6   129-192        24   11   4097-6144
    5    1    7,8         15   6   193-256        25   11   6145-8192
    6    2    9-12        16   7   257-384        26   12   8193-12288
    7    2    13-16       17   7   385-512        27   12   12289-16384
    8    3    17-24       18   8   513-768        28   13   16385-24576
    9    3    25-32       19   8   769-1024       29   13   24577-32768

Again, a very similar idea, for example if we have encountered code 4
According to the table above we need to read extra 1 bit, convert it to
a number and then add this number to the base distance:

    if the extra bit is set then
        distance = 5 + 1 = 6
    if the extra bit is not set then
        distance = 5 + 0 = 5

And this is it for the LZ77 portion.

----------
Now let's discuss Huffman coding algorithm. Why LZ77 is not enough to compress the data? As you can see we need
lots of bits to store length/distance information. What can be done to compress it further? What if somehow we
mapped both alphabets to shorter codes? An example
    suppose we have the next sequence to compress:
        ABCABC
        we can do this map:
        A - 00
        B - 01
        C - 11
    so instead of having ABCABC we use this sequence:
        00 01 11 00 01 11

    And when I say ABCABC in reality it means according to ASCII:
        A - 65
        B - 66
        C - 67
    or switching to bits it will look like this:
        01000001 01000010 01000011 01000001 01000010 01000011

    It's much much shorter when we use that map! As they say feel the difference!

Alright, this is the idea in general, replace long sequence of bits with their shorter versions. But there is a catch:
if we were to have a static map like with LZ77 alphabets it would not be that efficient. The next step is to
find the most frequently used symbols in the data we want to compress and assign the shortest codes to them.
It makes sense, consider this string for example:
        AAAABAAAACCAAAADAAAAEAAAA
        then if we assign 00 to A
        A - 00
        B - 01
        C - 100
        D - 101
        E - 111

        it will be the optimal map, just imagine if we assigned 111 to the A symbol...

So Huffman in 1951 while being a student came up with what later will be called in his honour - Huffman coding.
I don't want to go in full detail, as the theory that hides behind this simple facade is not small.
You can check wikipedia to see that it's quite complex. What is important for us is that there is some algorithm
that when given a sequence of bytes produces a map of much shorter codes:
    input AAAABAAAACCAAAADAAAAEAAAA
    Huffman(AAAABAAAACAAAADAAAAEAAAA) = our map
        A - 00
        B - 01
        C - 100
        D - 101
        E - 111
And it uses the number of occurrences to produce that map, it doesn't need the actual symbols:
    input AAAABAAAACCAAAADAAAAEAAAA
    from input we can see that the letter A happens 20 times in the data, B only once and so on
    20,1,1,1,1
    so the Huffman function when given this information will produce the next map:
     - 00
     - 01
     - 100
     - 101
     - 111
    and if we agree to pass the frequencies in the natural order i.e. a, b, c ... z, 0, 1, 2, ... 9
    then we know how to reconstruct the map, i.e. to which symbol the sequence needs to be assigned:
         A - 00
         B - 01
         C - 100
         D - 101
         E - 111
One more important thing to note, that each code in the Huffman codes map is perfectly unique for our bit
by bit reading. Let me illustrate:

    000000000100000000100000000001010000000011100000000
    <------------<--------------------------<----------

    step 1:
            we read 0 from the stream
                checking the map, the first code starts with 00, so there is no match
    step 2:
            we read 0 from the stream
                so now we have read 00 and there IS the corresponding sequence in our map
            so we output A and reset the read sequence
    step 3: we read 0 again, no matches, we move on to the next bit
    ...
    and so on
    until all the bits are read
Ok, what important for us at this point is that we don't need the actual symbols to produce the Huffman map.
Now let's go back to where we left it with the Dynamic Huffman BLOCK:
we said there that the BLOCK's DATA is:

    DATA begins with the following three integers:
        HLIT (5 bits) - number of literal/length codes - 257
        HDIST (5 bits) - number of distance codes - 1
        HCLEN (4 bits) - number of code length codes - 4

as you can see to get the actual number for literal/length codes we need to add 257 to it and so on:
    nlit = HLIT + 257

    ndist = HDIST + 1
    nclen = HCLEN + 4

while nlit and ndist are more or less clear (we will cover them in more details anyway), the HCLEN part
needs further explanation:
    the information needed to generate Huffman bit sequences for the length/distance (LIT) and
    for the distances (DIST) itself can occupy space, so the idea is to use one more Huffman map this time just to
    compress that information! Yeah, I hear you, it's quite overwhelming! This is why
    it can be so confusing. It took me some time to figure it out.

Ok, and to make things even more confusing this extra Huffman map uses its own variation of
LZ77, it's sort of a similar idea, but not exactly the same.

So to compress HLIT and HDIST there is this dictionary:

    16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15

Why in this order? It's based on the assumption that this particular order reflects that 16 is the
most frequently occurring number, and 15 is the least. It's just an assumption, but as we go
go forward it will become clear that it's just a common sense.

Anyway, it's time to show you the meaning behind these numbers:

    0 - 15: Represent code lengths of 0 - 15
    16: Copy the previous code length 3 - 6 times.
        The next 2 bits indicate repeat length
        (0 = 3, ... , 3 = 6)
        Example:
            Codes
                8, 16 (+2 bits 11), 16 (+2 bits 10)
            will expand to
                12 code lengths of 8 (1 + 6 + 5)
    17: Repeat a code length of 0 for 3 - 10 times.
        (3 bits of length)
    18: Repeat a code length of 0 for 11 - 138 times
        (7 bits of length)

I'll leave it here for a moment. Remember I told you that the whole Huffman coding is to
produce unique bit sequences and associate the shortest sequences with the most frequently
occurring symbols? So what does "code length" mean? It means the lengths of those unique
Huffman code sequences. An example:
    if we were given this information:
        code length          number of occurrences
        2                     1
        3                     3
    it would mean the next map:
        00                    1
        -----------------------
        010
        011                   3
        111

In other words our ultimate goal is to produce Huffman unique sequences and map these
sequences to the symbols we want to compress:
    so when we have
        the number of bits in a Huffman bit sequence
            +
        how many sequences of this given length
    we can generate the actual sequences!

As nn example let's have code length 4 and the number of occurrences 7:

        1000 \
        1001  \
        1010   \
        1011    7 sequences in total, and note that the next sequence = previous sequence + 1
        1100   /
        1101  /
        1111 /

If we had code length 4 and the number of occurrences 5 it would become:

        1000 \
        1001  \
        1010    5 sequences in total
        1011  /
        1100 /

The starting bit sequence for the given code length is a problem of course. How do we know that it begins
with 1000? Yeah while it's a bit tricky, eventually it's not that hard. But more about it later. What
important for us at the moment, that we sort of have an idea how to produce Huffman bit sequences!

Ok, back on track. We have these 3 integers:

    nlit = HLIT + 257
    ndist = HDIST + 1
    nclen = HCLEN + 4

What now? We read nclen * 3 bits, pseudocode
    for i from 0 to nclen
        bits = read(3 bits)

And remember we had this order 16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15?
It's time to use it! As we read those 3 bit long sequences we move along the dictionary:

    3 bits  - 16     |
    3 bits  - 17     V
    3 bits  - 18     |
    ...              V
    so on            .
    ...              .
    nclen times      |

3 bits represent the length of the Huffman bit sequence for the number in the dictionary, example:
    16 - 0 // 16 is not used to compress LIT and DIST information
    17 - 2 // 17 will be coded using 2 bits
    18 - 3 // 18 will be coded using 3 bits
    0  - 3 // 0  will be coded using 3 bits
    ...

Now we can calculate how many <code length>s we have and apply our + 1 algorithm to produce
Huffman bit sequences. From our example above we know that 3 bits long sequence will happen at least 2 times,
as 18 and 0 both have 3.

Ok let's imagine that we have generated Huffman bit sequences for the

    16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15

dictionary, now how to associate Huffman bit sequences with the numbers from the dictionary?
You may think that we assign them in the same order, something like if we have next sequences:

    100
    101
    ...

then the map will be:

    100 - 18
    101 - 0

WRONG! We need to use the natural order! I.e.:

    101 - 0
    100 - 18

as 0 < 18.


Having this map now we can decode LIT and DIST information:
    So we read nlit Huffman codes next. The CLEN map will give us:
numbers from the dictionary:

    16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15

Using pseudo code:
    output = []
    for i from 0 to nlit
        Huffman code = read bits until the read sequence matches with CLEN_MAP
        n = CLEN_MAP[Huffman code]
        if n < 15
            codeLen = n
            output.append(n)
        if n == 16
            bits = read(2 bits)
            prevIndex = len(output) - 1

            for j from 0 to 3 + toInt(bits)
                if prevIndex + j < 0 {
                        prevValue = 0
                } else {
                    prevValue = output[prevIndex + j]
                }
                output.append(prevValue)
            }
        if n == 17
            bits = read(3 bits)
            for j from 0 to 3 + toInt(bits)
                output.append(0)
            }
        if n == 18
            bits = read(3 bits)
            for j from 0 to 11 + toInt(bits)
                output.append(0)
            }
so in the output we can have something like:
    [0, 0, 0, 0, 2, 3, 3, 3, ....

And remember how many times I've mentioned the natural order? It applies here as well:

for LIT we have 286 symbols max:
so if we produced next Huffman bit sequences:

    00
    100
    101
    111
    ...

then the final LIT map would be:

    00  - 4 note that the numbers 0, 1, 2, 3 had 0 in the output, so it means that these numbers are not present
    100 - 5
    101 - 6
    111 - 7
    ....

For DIST we do the same:
   outputForDist = []
       for i from 0 to ndist
           Huffman code = read bits until the read sequence matches with CLEN_MAP
           n = CLEN_MAP[Huffman code]
           if n < 15
               codeLen = n
               outputForDist.append(n)
           if n == 16
               bits = read(2 bits)
               prevIndex = len(outputForDist) - 1

               for j from 0 to 3 + toInt(bits)
                    if prevIndex + j < 0 {
                        prevValue = 0
                    } else {
                        prevValue = output[prevIndex + j]
                    }
                   outputForDist.append(prevValue)
               }
           if n == 17
               bits = read(3 bits)
               for j from 0 to 3 + toInt(bits)
                   outputForDist.append(0)
               }
           if n == 18
               bits = read(3 bits)
               for j from 0 to 11 + toInt(bits)
                   outputForDist.append(0)
               }

and in the outputForDist we can get an array like this:
   [0, 0, 0, 3, 3, ...

Again we know the number of distance codes:
    from 0 to 29
so if have generated next Huffman bit sequences:
    000
    001
    010
    ....
it will mean the next DIST map:
    000 - 3
    001 - 4
    ....
Having LIT and DIST maps we can decompress the data, pseudo code:
    decompressedData = [] // if it's the very first BLOCK
                          // if it's not the first BLOCK we should already have some data in
                          // the decompressedData array
    for {
        Huffman code for LIT = read bits until the read sequence matches with LIT_MAP
        nLit = LIT_MAP[Huffman code for LIT]
        if nLit <= 255 {
            decompressedData.append(nLit)
            continue
        }

        if nLit == 256 { // it's the stop token
            break        // we need to read the HEADER of the next BLOCK
        }
        // if  n > 256
        extraBitsForLen = read the number of extra bits according to nLit (see the table for LIT)
        len = baseLen for nLit + toInt(extraBitsForLen)

        //  now we read Huffman code for DIST
        Huffman code for DIST = read bits until the read sequence matches with DIST_MAP
        nDist = DIST_MAP[Huffman code for DIST]
        extraBitsForDist = read the number of extra bits according to nDist
        dist = baseDist for nDist + toInt(extraBitsForDist)
        for i from 0 to len {
            decompressedData.append(decompressedData[-dist + i]) // we go back dist times in
                                                                 // decompressedData array
        }
    }
Note also that the referenced string may overlap the current position;
for example, if the last 2 bytes decoded have values X and Y, a string reference with
    <length = 5, distance = 2>
we will need to add X,Y,X,Y,X to the output

I think the last piece that is left to be covered is the full algorithm of how
Huffman bit sequences are produced. If we know the starting sequence for the given code length
we simply add + 1, but the big question is how to get that starting sequence for a given code length.
Let's start with an example again:
    If we have the next data:

        code length        number of codes
        2                  1
        3                  2
        4                  7

    we begin with 00 the first sequence is kinda intuitive
    then we do next:
        nextSequence = (prevSequence + 1) << 1
    so for 3 we start with:
        00 + 1 = 01
        01 << 1 = 010 // it's the starting sequence for code length 3:
    we have 2 such sequences:
        010
        011
    ok time to go to the next sequence:
        011 + 1 = 100
        100 << 1 = 1000 // it's the starting sequence for code length 4:
    we have 7 such sequences:
        1000
        1001
        1010
        1011
        1100
        1101
        1111

So far so good, problems begin if we have gaps in code lengths:

        code length        number of codes
        3                  2
        5                  23

Ok for code length 3 it is still easy:
        000
        001
But how to jump to code length 5?
        001 + 1 = 010
        010 << 1 = 0100 // 4 but we need 5!
And the solution is:
        we continue shifting until we get needed code length
        0100 << 1 = 01000 // it's exactly what we need!
As we have 23 such codes we keep adding + 1 to the starting sequence:
        01000
        01001
        01010
        01011
        01100
        01101
        01110
        01111
        10000
        10001
        10010
        10011
        10100
        10101
        10110
        10111
        11001
        11010
        11011
        11100
        11101
        11110
        11111

This is the end for Dynamic Huffman BLOCK type

================================================================================================
The Fixed Huffman codes BLOCK is similar to Dynamic but the
LIT and DIST maps are hardcoded, so you can start decoding the data right away.

	 Lit Value             Bits            Codes
	 0-143                 8               00110000  - 10111111
	 144 - 255             9               110010000 - 111111111
	 256 - 279             7               0000000   - 0010111
	 280 - 287             8               11000000  - 11000111

So if the Huffman code is 00110000 it means 0, 00110001 - 1, and so on...

Distance codes 0-31 are represented by (fixed-length) 5-bit codes

================================================================================================
The final words. I want to share some of the links I used when was trying to understand PNG format:
http://www.libpng.org/
https://handmade.network/forums/articles/t/2363-implementing_a_basic_png_reader_the_handmade_way
https://commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art001
https://blog.za3k.com/understanding-gzip-2/
https://www.rfc-editor.org/rfc/rfc1951#section-3.1.1
