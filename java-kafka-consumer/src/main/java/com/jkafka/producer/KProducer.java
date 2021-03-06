package com.jkafka.producer;

import com.jkafka.Logger;
import java.util.Properties;
import java.util.UUID;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.serialization.StringSerializer;

public class KProducer {

  private final String[] brokers;
  private final String topic;

  public KProducer(String[] brokers, String topic) {
    this.brokers = brokers;
    this.topic = topic;
  }

  public String isDownVoted(String message) {
    String[] str = message.split("::");
    if (str.length < 2) {
      Logger.log("Not A Valid Vote");
      return "0099-0099-SKIP";
    }
    if ("down".equalsIgnoreCase(str[1])) {
      return str[0];
    }
    return "0099-0099-SKIP";
  }

  public void sendMessage(String message) {
    String candidate = isDownVoted(message);
    if (candidate.equalsIgnoreCase("0099-0099-SKIP")) {
      return;
    }
    Properties props = new Properties();
    String kafkaBrokers = String.join(",", brokers);
    props.setProperty(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG,
        kafkaBrokers);
    props.setProperty(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,
        StringSerializer.class.getName());
    props.setProperty(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG,
        StringSerializer.class.getName());
    props.setProperty(ProducerConfig.ACKS_CONFIG, "1");

    try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {
      String key = UUID.randomUUID().toString();
      ProducerRecord<String, String> record = new ProducerRecord<>(topic, key, candidate);
      producer.send(record);
      Logger.log("DownVoted For Candidate " + candidate);
    } catch (Exception ex) {
      Logger.log(ex.getMessage());
    }

  }
}
