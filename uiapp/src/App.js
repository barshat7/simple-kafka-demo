import './App.css';
import Button from './components/Button';
import VoteButtons from './components/VoteButtons';
import {useState, useEffect} from 'react';
import LeaderBoard from './components/LeaderBoard';

function App() {

  const BACKEND_URL = "http://localhost:3010";

  const [leader, setLeader] = useState("Not Defined Yet");

  const candidates = [
    {
      uniqueID: "Candidate-A",
      color: "blue"
    },
    {
      uniqueID: "Candidate-B",
      color: "green"
    },
    {
      uniqueID: "Candidate-C",
      color: "orange"
    }
  ]

  const upVote = async (uniqueID) => {
    const vote = {
      type: "up",
      uniqueID: uniqueID
    }
    await sendPost(vote);
  }

  const downVote = async (uniqueID) => {
    const vote = {
      type: "down",
      uniqueID: uniqueID
    }

    await sendPost(vote);
  }

  useEffect(() => {
    const sse = new EventSource(`${BACKEND_URL}/sse`,
      { withCredentials: false });
    function getRealtimeData(data) {
      console.log(`Leader Now ${data}`)
      setLeader(data)
    }
    sse.onmessage = e => getRealtimeData(e.data);
    sse.onerror = () => {
      // error log here 
      
      sse.close();
    }
    return () => {
      sse.close();
    };
  }, []);

  const sendPost = async (vote) => {
    const response = await fetch(`${BACKEND_URL}/vote`, {
      method: "POST",
      headers: {
        "Content-type": "Application/json"
      },
      body: JSON.stringify(vote)
    })
    const data = await response.json();
    console.log(`Voted Successfully ${JSON.stringify(data)}`)
  }

  return (
    <div className='body-class'>
      <h2 className = 'header'>Online Voting</h2>
      <VoteButtons type = {'upvote'}candidates = {candidates} onVote = {upVote} />
      <VoteButtons type = {'downvote'} candidates = {candidates} onVote = {downVote} />
      <LeaderBoard leaderName = {leader} />
    </div>
  );
}

export default App;
