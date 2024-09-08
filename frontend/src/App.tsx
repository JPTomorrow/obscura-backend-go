import "./App.css";
import YouTube from "react-youtube";
import {
  IconCopy,
  IconThumbUp,
  IconThumbDown,
  IconBrandDiscord,
  IconInfoCircle,
  IconPlayerSkipForward,
  IconX,
  IconSettings,
} from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { useCookies } from "react-cookie";

interface YTVideo {
  id: number;
  title: string;
  description: string;
  video_tag: string;
  upvotes: number;
  downvotes: number;
}

function App() {
  const [vid, setVid] = useState<YTVideo>({
    id: -1,
    title: "rick roll",
    description: "rick roll",
    video_tag: "dQw4w9WgXcQ",
    upvotes: 0,
    downvotes: 0,
  });
  const [isTopBarShowing, setIsTopBarShowing] = useState(true);
  const [isShowingAboutPanel, setIsShowingAboutPanel] = useState(false);
  const [cookies, setCookie] = useCookies(["name"]);

  let topBarInterval: number;
  const hideTopBar = (delay: boolean = true) => {
    if (!delay) {
      setIsTopBarShowing(false);
      return;
    }

    setIsTopBarShowing(true);
    clearTimeout(topBarInterval);

    topBarInterval = setTimeout(() => {
      setIsTopBarShowing(false);
    }, 5000);
  };

  const nextVideo = async () => {
    let resp = await fetch("http://localhost:8080/next-vid");
    let data = await resp.json();
    setVid(data);
  };

  useEffect(() => {
    nextVideo();
    hideTopBar();
  }, []);

  useEffect(() => {
    console.log(vid);
  }, [vid]);

  // const increaseVolume = (e: any) => {
  //   setTimeout(
  //     () => {
  //       e.target.setVolume(30);
  //       e.target.playVideo();
  //     },
  //     1000,
  //     e
  //   );
  // };

  const copyUrlToClipbard = () => {
    navigator.clipboard.writeText(
      "https://www.youtube.com/watch?v=" + vid.video_tag
    );
  };
  const handleThumbsUp = () => {};
  const handleThumbsDown = () => {};

  const opts = {
    playerVars: {
      // https://developers.google.com/youtube/player_parameters
      autoplay: 1,
      mute: 1,
      rel: 0,
      controls: 1,
      origin: "http://localhost:5173",
    },
  };
  return (
    <div className="black-bg">
      <YouTube
        onMouseMove={() => {
          console.log("mouse moved");
          hideTopBar();
        }}
        key={vid.id}
        videoId={vid.video_tag}
        onEnd={nextVideo}
        onPause={hideTopBar}
        onPlay={() => hideTopBar(false)}
        opts={opts}
        className="w-full h-full"
      />
      <div
        className={`fixed flex top-0 w-screen justify-center transition-transform ${
          isTopBarShowing ? "" : "-translate-y-24"
        }`}
      >
        <div className="bg-black/50 rounded-b-xl text-white h-fit w-fit flex flex-col gap-2 p-5">
          <div className="flex gap-2">
            <button
              onClick={copyUrlToClipbard}
              className="btn-base bg-gray-500"
            >
              <IconCopy size={28} className="text-white" />
            </button>
            <button onClick={handleThumbsUp} className="btn-base bg-green-700">
              <IconThumbUp size={28} className="text-white bg-green-700" />
            </button>
            <button onClick={handleThumbsDown} className="btn-base bg-red-800 ">
              <IconThumbDown size={28} className="text-white" />
            </button>
            <button
              onClick={() => nextVideo()}
              className="btn-base bg-red-800 "
            >
              <IconPlayerSkipForward size={28} className="text-white" />
            </button>
            <a
              href="https://discord.gg/SmTk8AedgG"
              target="_blank"
              className="btn-base bg-purple-800 "
            >
              <IconBrandDiscord size={28} className="text-white" />
            </a>
            <button
              onClick={() => setIsShowingAboutPanel(true)}
              className="btn-base bg-gray-600"
            >
              <IconInfoCircle size={28} className="text-white" />
            </button>
            <button
              onClick={() => setIsShowingAboutPanel(true)}
              className="btn-base bg-gray-600"
            >
              <IconSettings size={28} className="text-white" />
            </button>
          </div>
        </div>
      </div>
      <motion.div
        initial={{
          opacity: 1,
        }}
        animate={{
          opacity: 0,
          display: "none",
        }}
        transition={{
          delay: 3,
        }}
        className="fixed flex z-50 items-center justify-center bg-black top-0 left-0 w-screen h-screen"
      >
        <motion.div
          initial={{
            scale: 0,
          }}
          animate={{
            scale: [0.0, 1.1, 1.0],
          }}
          transition={{
            delay: 0.5,
            duration: 0.8,
          }}
          className="carnival-font w-fit flex flex-col gap-5 text-center"
        >
          <h1 className="text-[128pt] font-bold uppercase">Obscura</h1>
          <p className="text-7xl">Boo! That scared you, right?</p>
        </motion.div>
      </motion.div>
      {/* About Panel */}
      <Panel isShowing={isShowingAboutPanel}>
        <button className="float-right bg-transparent">
          <IconX size={42} onClick={() => setIsShowingAboutPanel(false)} />
        </button>
        <h1 className="text-6xl font-bold">{vid.title}</h1>
        <p className="text-3xl">{vid.description}</p>
      </Panel>
      {/* Settings Panel */}
    </div>
  );
}

const Panel = ({
  children,
  className,
  isShowing,
}: {
  children?: React.ReactNode;
  className?: string;
  isShowing: boolean;
}) => {
  const variants = {
    hidden: {
      opacity: 0,
      display: "none",
    },
    visible: {
      opacity: 1,
      display: "flex",
    },
  };
  return (
    <motion.div
      variants={variants}
      animate={isShowing ? "visible" : "hidden"}
      className={`fixed flex justify-center items-center w-screen h-screen top-0 left-0 ${className}`}
    >
      <div className="w-1/2 h-fit flex flex-col gap-5 bg-gray-800 p-10 rounded-3xl">
        <div className="overflow-y-auto max h-[200px]">{children}</div>
      </div>
    </motion.div>
  );
};

export default App;
