package com.jkafka;

import com.jkafka.client.KafkaClient;
import com.jkafka.producer.KProducer;
import java.util.Scanner;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;


public class Main {

  public static void main(String... args) {
    String[] brokers = {"localhost:9092"};
    String topic = "vote_topic";
    String producerTopic = "down_vote_topic";

    KProducer producer = new KProducer(brokers, producerTopic);

    ExecutorService executorService = Executors.newFixedThreadPool(2);
    executorService.submit(() -> {
      var kafkaClient = new KafkaClient(brokers, topic, "java-consumer");
      kafkaClient.consume(producer::sendMessage);
    });

    Scanner scanner = new Scanner(System.in);
    Logger.log("Press Any Key To Stop");
    scanner.next();
  }
}
