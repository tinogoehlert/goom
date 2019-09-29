[Goom Documentation](README.md) > [MUS File Format](MUSIC.md)

Links
-----
[1]: http://www.shikadi.net/moddingwiki/MUS_Format            MUS Format (ModdingWiki)
[2]: https://www.doomworld.com/idgames/docs/editing/mus_form  MUS Format Description (by Vladimir Arnost)

MUS File Format
===============

A MUS file is defined as follows:

    type MusData struct {
        []byte    Id           // The 4-byte identifer must be "MUS\x1a"
        uint16    scoreLen     // Number of bytes of score data
        uint16    scoreStart   // Position of the first score data byte
        uint16    channels	   // Number of primary channels (excl. percussion channel 15)
        uint16    secChannels  // Number of secondary channels
        uint16    numInstr     // Number of instruments
        uint16    dummy        // Separator
        []uint16  instruments  // List of used instruments
        []bytes   scores       // Music events
    }

Instruments
-----------

    InstrumentNumber     Description
      0 - 127            standard MIDI instruments
    135 - 181            standard MIDI percussions (notes 35 - 81)

Mus Events
----------

    ╓──7─┬──6─┬──5─┬──4─┬──3─┬──2─┬──1─┬──0─╖
    ║Last│   MusEvent   │   ChannelNumber   ║
    ╙────┴────┴────┴────┴────┴────┴────┴────╜

`ChannelNumber`

The `ChannelNumber` specifies, which audio settings to use.
Each channel defines a setting as (`InstrumentNumber`, `Panning`, `Volume`).
The regular channels are `0` to `14`, channel `15` is used only for percussions.

`MusEvent`

    // MusEvent types.
    const (
        RelaseNote       = 0
        PlayNote         = 1
        PitchWheel       = 2
        SystemEvent      = 3
        ChangeController = 4
        Unknown5         = 5
        ScoreEnd         = 6
        Unknown7         = 7
    )

`Last`

Indicates the last event in a group of concurrent events.
The last event is followed by a duration in `Ticks` to wait
before proceeding with the next group.

    Tick Length             Game
    1s / 140  ~=  7.1ms     Doom I, II, and Heretic
    1s / 70   ~= 14.2ms     Raptor

`Tick` Arithmetic

*Work in Progress! The followign may not be correct or complete!*

> *Source*: [MUS Format Description by Vladimir Arnost][2]
> 
> 1. time = 0
> 2. READ a byte
> 3. time = time * 128 + byte AND 127
> 4. IF (byte AND 128) GO TO 2
> 5. RETURN time
> 
> The time info is a series of 7-bit chunks.
> The bit #7 is set until the last byte whose bit 7 is zero.
> This scheme allows small numbers occupy less space than large ones.

Example:

    In  Last Evt Ch  Effect                Values                             Last Effect
    80     1 000  0  RelaseNote event (0); Delay = 0                          1 -> delay follows
    10     0 001  0  RelaseNote data:      Note = 1 (Stop playing note 1)     none
    82     1 000  2  Add delay;            Delay = Delay * 128 + 2 = 2        1 -> delay follows
    05     0 000  5  Complete delay;       Delay = Delay * 128 + 5 = 261      0 -> delay complete.
