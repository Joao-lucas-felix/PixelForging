package com.pixel_forging.pixel_forging.controllers;

import com.pixel_forging.pixel_forging.dtos.StatusDto;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("/wake")
public class GreentingController {
    @GetMapping
    public ResponseEntity<StatusDto> wake() {
        return  ResponseEntity.ok(new StatusDto("UP!"));
    }
}
