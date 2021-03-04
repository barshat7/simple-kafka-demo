import './App.css';
import Button from './components/Button';
import TextInput from './components/TextInput';
import {useState, useEffect} from 'react';

function App() {

  const BACKEND_URL = "http://localhost:3010";

  const [eventAgree, setEventAgree] = useState("");
  const [eventDisAgree, setEventDisAgree] = useState("");

  const createEventAgree = async () => {
    let _eventName = eventAgree;
    if (!_eventName) {
      _eventName = 'Default Event:- ' + Date() ;
    }
    const event = {
      type: 'agree',
      value: _eventName
    }
    await sendPost(event);
  }
  
  const createEventDisAgree = async () => {
    let _eventName = eventDisAgree;
    if (!_eventName) {
      _eventName = 'Default Event:- ' + Date() ;
    }
    const event = {
      type: 'disagree',
      value: _eventName
    }
    await sendPost(event);
  }

  const sendPost = async (body) => {
    const res = await fetch(`${BACKEND_URL}/events`, {
      method: "POST",
      headers: {
        "Content-type": "application/json"
      },
      body: JSON.stringify(body)
    });
    const data = await res.json();
    console.log(`Event Created ${JSON.stringify(data)}`);
  }

  const onEventAgreeSet = (e) => {
    setEventAgree(e.target.value);
  }

  const onEventDisAgreeSet = (e) => {
    setEventDisAgree(e.target.value);
  }

  return (
    <div className='body-class'>
      <h2 className = 'header'>Kafka Pipeline</h2>
      <TextInput onChange = {onEventAgreeSet} />
      <Button text='Agree' onClick = {createEventAgree}/>
      <TextInput onChange = {onEventDisAgreeSet} />
      <Button text='Disagree' onClick = {createEventDisAgree}/>
    </div>
  );
}

export default App;
