SELECT add_continuous_aggregate_policy('stats_threat_by_src',
                                       start_offset => INTERVAL '1d',
                                       end_offset   => INTERVAL '1h',
                                       schedule_interval => INTERVAL '10 minutes',
                                       if_not_exists => true);

SELECT add_continuous_aggregate_policy('stats_threat_by_dst',
                                       start_offset => INTERVAL '1d',
                                       end_offset   => INTERVAL '1h',
                                       schedule_interval => INTERVAL '10 minutes',
                                       if_not_exists => true);

SELECT add_continuous_aggregate_policy('stats_threat_by_severity',
                                       start_offset => INTERVAL '1d',
                                       end_offset   => INTERVAL '1h',
                                       schedule_interval => INTERVAL '10 minutes',
                                       if_not_exists => true);

SELECT add_continuous_aggregate_policy('stats_threat_by_phase',
                                       start_offset => INTERVAL '1d',
                                       end_offset   => INTERVAL '1h',
                                       schedule_interval => INTERVAL '10 minutes',
                                       if_not_exists => true);