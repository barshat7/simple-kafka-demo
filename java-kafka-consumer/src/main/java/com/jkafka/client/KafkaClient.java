package com.jkafka.client;


import com.jkafka.Logger;
import java.time.Duration;
import java.util.Collections;
import java.util.Properties;
import java.util.function.Consumer;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;

public class KafkaClient {

  private final String[] brokers;
  private final String topic;
  private final String groupID;

  public KafkaClient(String[] brokers, String topic, String groupID) {
    this.brokers = brokers;
    this.topic = topic;
    this.groupID = groupID;
  }

  public void consume(Consumer<String> consumer) {
    Logger.log("Starting Kafka Client...");
    Properties props = new Properties();
    String kafkaBrokers = String.join(",", brokers);
    props.setProperty(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG,
        kafkaBrokers);
    // user-defined string to identify the consumer group
    props.put(ConsumerConfig.GROUP_ID_CONFIG, groupID);
    // Enable auto offset commit
    props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, "true");
    props.put(ConsumerConfig.AUTO_COMMIT_INTERVAL_MS_CONFIG, "1000");
    props.setProperty(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG,
        StringDeserializer.class.getName());
    props.setProperty(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG,
        StringDeserializer.class.getName());
    try (KafkaConsumer<String, String> kafkaConsumer = new KafkaConsumer<>(props)) {
      kafkaConsumer.subscribe(Collections.singleton(topic));
      while (true) {
        try {
          ConsumerRecords<String, String> records = kafkaConsumer.poll(Duration.ofMillis(10000));
          records.forEach(record -> consumer.accept(record.value()));
        } catch (Exception ex) {
          Logger.log(ex.getMessage());
          break;
        }
      }
    } catch (Exception ex) {
      throw new IllegalStateException(ex);
    }
  }
}
