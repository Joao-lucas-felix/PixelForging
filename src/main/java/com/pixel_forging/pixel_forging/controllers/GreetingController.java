package com.pixel_forging.pixel_forging.controllers;

import com.pixel_forging.pixel_forging.dtos.StatusDto;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/wake")
public class GreetingController {
    private static final Log log = LogFactory.getLog(GreetingController.class);

    @GetMapping
    public ResponseEntity<StatusDto> wake() {
        log.info("The server is WakeUp!");
        return  ResponseEntity.ok(new StatusDto("UP!"));
    }
}
