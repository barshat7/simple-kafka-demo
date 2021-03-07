const {Kafka} = require('kafkajs');

const main = async () => {
    await createDestinationTopic();
    await upvoteProcessor();
    //await downVoteProcessor();
}

const destinationTopic = "leaderboard_topic";

let candidates = new Map();

const createDestinationTopic = async () => {
    const  kafka = new Kafka({
        "clientId": "myapp",
        "brokers": ["127.0.0.1:9092"]
      });
      const admin = kafka.admin();
      console.log('Connecting...')
      await admin.connect();
      console.log('Connected');
    
      await admin.createTopics({
        topics: [{
          topic:  destinationTopic
        }]
      });
      console.log('Topic Created Successfully');
      await admin.disconnect();
}

const sendToDestinationTopic = async (message) => {
    const  kafka = new Kafka({
        "clientId": "leaderboardservice",
        "brokers": ["127.0.0.1:9092"]
      });
      const producer = kafka.producer();
      await producer.connect();
      const response =  await producer.send({
        topic: destinationTopic,
        messages: [
          {
            value: message
          }
        ]
      })
      console.log('Message sent to leaderboard topic ');
      await producer.disconnect();
}


const upvoteProcessor = async () => {
    const  kafka = new Kafka({
        "clientId": "leaderboardprocessor",
        "brokers": ["127.0.0.1:9092"]
      });
      const consumer = kafka.consumer({
        "groupId": "test_1" // We are defining a consumer group here.. name is test
      });
      console.log('Connecting...')
      await consumer.connect();
      console.log('Connected');
      
      await consumer.subscribe({
        topic: "up_vote_topic", // This consumer will subscribe to users topic
        fromBeginning: false
      });
      await consumer.subscribe({
          topic: "down_vote_topic",
          fromBeginning: false
      })
      console.log("Subscribed");
    
      await consumer.run({
        eachMessage: async ({ partition, message, topic }) => {
          console.log(`UpVote  Message: ${message.value.toString()} On Topic ${topic} On Partition ${partition}`);
          writeToCandidate(message.value.toString(), topic);
          await sendMaxMin();
        }
      })
}

const writeToCandidate = (candidateName, topic) => {
    let count;
    if (topic == 'up_vote_topic') {
        count = Number(1);
    } else {
        count = Number(-1);
    }
    if (!candidates[candidateName]) {
        candidates.set(candidateName, count);
    } else {
        candidates.set(candidateName, candidate[candidateName] +  count);
    }
}

const sendMaxMin = async () => {
    console.log(JSON.stringify(candidates));
    let maxCandidateName;
    let minCandidateName;
    let max = Number(-100);
    let min = Number(999999);
    for (const [key, value] of candidates.entries()){
        console.log("Found Key " +key);
        if (value > max) {
            max = value;
            maxCandidateName = key;
        }
    }
    for (const [key, value] of candidates.entries()){
        if (value < min && key !== maxCandidateName) {
            min = value;
            minCandidateName = key;
        }
    }
    await sendToDestinationTopic(maxCandidateName);
}

main().catch(async error => {
    console.error(error)
    try {
      await consumer.disconnect()
    } catch (e) {
      console.error('Failed to gracefully disconnect consumer', e)
    }
    process.exit(1)
  })