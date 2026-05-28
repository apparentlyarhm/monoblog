---
title: "The Cake is Javascript: You can over-engineer anything!"
category: "Tech"
date: "2026-03-18"
slug: "the-cake-is-javascript"
tags: "javascript, typescript, raf, React, Portal, producer/consumer"
description: "Re-creating the Portal end credits animation with a Producer/Consumer architecture"
image: "/cover-cake-is-js.png"
---

# I. Preface

I'm not gonna lie to you, [Portal](https://store.steampowered.com/app/400/Portal/) has been around since forever and yet, It was only recently that I decided to play it through; Absolutely adored the both the games and given I was already a fan of the Half Life series, the universe inter-twining and references of one in the other was the perfect cherry on top for me and it was no surprise that I loved the games as much as I did. While I'm done waiting for Half Life 3, There's something else that caught my eye (or ears, debatable) in the first game - the End Credit song - [<i>"Still Alive."</i>](https://www.youtube.com/watch?v=Y6ljFaKRTrI&list=RDY6ljFaKRTrI&start_radio=1)   



The animation in this, is so elegant, yet, simple. *"I should make this for me website!"*, I thought. <i>"Looks simple enough"</i> and to be honest, it was, and yet, to me, it was not entirely clear as to how I would approach something like this. I had to reasearch.

Shortly after however, I came across [this repo](https://github.com/errorer/Portal_StillAlive_Python). This was perfect, since I now have a reference of the timing of the lyrics and corresponding action it does. One small thing tho, how do we translate this to JS on the browser?

# II. Building the animation

## Animations have a tiny problem in general that has been solved

At first, I thought I could probably get away with some cursed combination of `setTimeout`s and chained delays. And to be fair, for tiny demos? It *kind of* works.

This is the case until you realize timing on the browser is... unpredictable.

Apart from refresh rates being different nowadays, there are a lot of actions that can trigger the timing to behave predicatably - like moving out of tabs, or perhaps browser frozen because of other tabs. There are just too many 

Remember pokemon? NDS? GBA? Try emulating it using your modern, state-of-the-art computer. Now, remove the FPS cap. The game PROPORTIONALLY speeds up. The logic was essentially:

**"1 frame = 1 unit of game progression."**

This was fine back when hardware was predictable.

I realized I was dangerously close to doing the exact same thing here.

If I tied lyric progression directly to frames, then the animation would literally run faster on higher refresh-rate monitors. A 144Hz display would process updates more frequently than a 60Hz one.

That was essentially my introduction to `requestAnimationFrame` - but dont get me wrong - this *isn't* the solution to our problem, it just makes everything elegant and easier.

## So what is requestAnimationFrame?

As I understand it, its just the browser saying - 

<i>"Hey, I'm about to repaint the screen. If you want to update anything visual, now would be a good time."</i>

You give it a callback, and the browser calls that function right before the next repaint. On a 60Hz monitor, this usually means around 60 calls per second. On a 144Hz display, closer to 144.

*Why is raF not the complete solution here then?*

Just imagine this:

```js
function step() {
    position += 1;
    requestAnimationFrame(step);
}
```

This looks innocent enough until you realize that on a higher refresh-rate monitor, this callback executes more frequently. Which means `position` increments faster. Which means the animation itself speeds up.

This is actually very similar to the Pokemon problem we discussed.

The important realization for me was this:

<b>Frames are not time.</b>

Frames are just opportunities to *check* time.

That is why the timestamp argument provided by `requestAnimationFrame` ended up becoming the foundation of the entire timing system.

## Wait hold on - timing *system?*

So far what we do have is
- ***Timing data of the lyrics*** - when each lyric line / ASCII art is triggered, along with internvals of each letter.
- **raF** - a way to call a function automatically based on framerate.

Right now, all we really have is:

- a list of events
- and a function that keeps getting called.

Those two things alone do not magically create synchronization. At some point, we need to answer a much more important question:

**"Given the current point in time, which events should have happened by now?"**

The browser calling our function every frame is nice, but the callback itself has no idea:

- which lyric should currently be visible
- whether ASCII art should be drawn yet
- whether music has started
- whether multiple events were missed because of lag
- or whether the tab froze for half a second and suddenly resumed

We need something sitting in the middle constantly watching the clock and dispatching events at the correct moments.

Essentially, what we need is a scheduler. It could be something like:

```js
useEffect(() => {
    let currentEventIndex = 0;
    function step(timestamp) {
        // figure out how much time has passed 
        const elapsedTime = timestamp - startTime;

        /* 
            Keep processing events until we catch up to 
            where the timeline says we should be. 
        */
        while ( 
            currentEventIndex < timelineEvents.length && 
            elapsedTime >= timelineEvents[currentEventIndex].time) {
                const event = timelineEvents[currentEventIndex];
                // this `event` data structure is completly arbitrary and can be achieved in multiple ways.
                switch (event.mode) { 
                    case 'LYRIC': // queue lyrics drawing
                    case 'DRAW_ART': // queue drawing
                    case 'START_MUSIC': // just trigger the music
                } 
        
            currentEventIndex++;
       }
       requestAnimationFrame(step);
    }
    requestAnimationFrame(step);
}, [])
```
where `timelineEvents` is an array of all events that the song needs: 

```ts
export interface LyricLine {
    words: string | number; // The lyric line to display. if its a number then its ascii art index
    time: number; // Time in centiseconds when this line should appear
    interval: number; // Duration to display the line (optional)
    mode: 'LYRIC_NEWLINE' | 'START_MUSIC' | 'DRAW_ART' | 'CLEAR_LYRICS' | 'LYRIC_NONEWLINE' | 'END'; // Mode of the line
}

export const timelineEvents: LyricLine[] = [
    // { words, time, interval, mode }
    { words: "Forms FORM-29827281-12:", time: 0, interval: -1, mode: 'LYRIC_NEWLINE' },
    { words: "Test Assessment Report", time: 200, interval: -1, mode: 'LYRIC_NEWLINE' },
    { words: "\x00\x00\x00\x00\x00\x00\x00", time: 400, interval: -1, mode: 'LYRIC_NEWLINE' },
    { words: "", time: 710, interval: 0, mode: 'START_MUSIC' },
    { words: "This was a triumph.", time: 730, interval: 2, mode: 'LYRIC_NEWLINE' },
    ...
]
```

The responsibility of this main scheduler is therefore to:
- continuously check the current time
- compare it against the timeline
- dispatch work when needed

Since I have omitted the "credits" section of the original animation, the "work" here is to either draw ASCII or type lyrics- all while being in sync. neat.

What is even better here is the fact that this system, by default, is "pause-ready". that means we dont really need to do much to support animation pausing with clean resumes.

## The "consumer" - useEffects for both typing and drawing

The title here is pretty self explanatory, ommitting the implementation details of 2 functions - `typeLyrics` and `drawAscii` for now, we need to plug these in into 2 separate useEffects that will act as a consumer for all the "tasks" the scheduler dispatches: 

```js
useEffect(() => {

    /* 
    Only continue if: 
    - there is work available 
    - we are not already processing another task 
    */

    if (!isProcessing && taskQueue.length > 0) {

      // Peek at the first item (FIFO - First In First Out)
      const task = taskQueue[0];

      (async () => {
        setIsProcessing(true); // lock consumer

        /* 
            Figure out typing speed. 
            Some events specify explicit intervals. 
            Others dynamically calculate intervals based 
            on the next timeline event. 
        */

        await typeLyrics(task.words, interval); // our actual worker

        setTaskQueue(currentQueue => currentQueue.slice(1)); // removed from q
        setIsProcessing(false); // unlock for next
      })();
    }
  }, [taskQueue, isProcessing]);
```

the dep array here is interesting. Whenever the queue changes, or whenever processing state changes, we re-evaluate whether more work needs to be done. thats it.


## This is similar to producer-consumer system in spirit, not in practice

Looking back, the architecture ended up feeling surprisingly similar to a tiny producer/consumer system.

The scheduler continuously scans the timeline and produces work:

- start typing this lyric
- draw this ASCII art
- start music
- clear the screen

Meanwhile, independent `useEffects` consume those tasks and execute the actual visual behavior.

I definitely wasn't thinking in those terms while building it, but in hindsight, separating "detecting when something should happen" from "performing the thing itself" made the whole system dramatically easier to reason about. This isnt *"Today we will implement a producer-consumer architecture in React™"* anyway.

The workers here don't add too much to this discussion, so for those who are interested, can check out the code in the [repo](https://github.com/apparentlyarhm/notaportfolio/commit/459b778d6432d83b7358fa53ec4fef3c88f64818). Keep in mind this is the initial version of the engine with the basic bells and whistles we discussed so far. Since then, it has seen some changes, but mostly minor.

## A quick napkin math to solidify all this even further

Suppose we have this timeline event:

```
{ 
    words: "This was a triumph.",
    time: 730
}
```
Remember, the timeline uses centiseconds. So this event should trigger at:


$$730\text{ centiseconds} = 7300\text{ ms} = 7.3\text{ seconds}$$

Now imagine the animation starts at browser timestamp:

$$ 50000ms $$

And let's say the user is on a 144Hz monitor. That means `requestAnimationFrame` gets called roughly every:

$$\frac{1000}{144} \approx 6.94\text{ ms/frame}$$

So eventually, after enough frames, our callback timestamps will look something like:

```
50000
50006.94
50013.88
50020.82
...
```

The scheduler continuously calculates:

$$ elapsedTime = timestamp - startTime$$

So what we really want to know is:

"At approximately which frame does `elapsedTime` become greater than 7300ms?"

That becomes:

$$N \cdot 6.94 \geq 7300$$

Solving for N:

$$N \geq \frac{7300}{6.94} \approx 1051.87$$

Meaning somewhere around frame 1052, the scheduler notices:

$$ elapsedTime >= event.time $$

and dispatches the lyric task.

Of course, in reality browsers are not perfectly precise clocks. Frames can jitter slightly, callbacks can be delayed, and the browser may skip frames under load.

# III. Closing Thoughts

I really should stop adding this section, as there really isnt much to say as the outro. Still, you can find the animation live [here](https://space.arhm.dev/still-alive). Hopefully this was an interesting usecase of various tools and softwares that are available to us as developers.